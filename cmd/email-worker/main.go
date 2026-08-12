package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/ticketing-service/config"
	"github.com/mbeka02/ticketing-service/internal/broker/rabbitmq"
	"github.com/mbeka02/ticketing-service/internal/email"
	resendmailer "github.com/mbeka02/ticketing-service/internal/mailer/resend"
	"github.com/mbeka02/ticketing-service/internal/outbox"
	"github.com/mbeka02/ticketing-service/internal/postgres"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("load configuration", zap.Error(err))
	}
	if err := cfg.ValidateEmailWorker(); err != nil {
		logger.Fatal("invalid email worker configuration", zap.Error(err))
	}
	if err := logger.Init(cfg.ServerEnv); err != nil {
		logger.Fatal("initialize logger", zap.Error(err))
	}
	defer logger.Sync()

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal("connect to RabbitMQ", zap.Error(err))
	}
	defer conn.Close()
	if err := (rabbitmq.Topology{}).Setup(conn); err != nil {
		logger.Fatal("declare RabbitMQ topology", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	store, err := postgres.NewStore(ctx, cfg.GetDatabaseConfig())
	if err != nil {
		logger.Fatal("connect to PostgreSQL", zap.Error(err))
	}
	defer store.Close()

	sender, err := resendmailer.New(cfg.ResendAPIKey, cfg.EmailFrom)
	if err != nil {
		logger.Fatal("create email sender", zap.Error(err))
	}
	publisher := rabbitmq.Publisher{Conn: conn, ConfirmTimeout: 5 * time.Second}
	owner := relayOwner()
	relay := outbox.Relay{
		Store: postgres.NewOutboxRepository(store), Publisher: publisher, Owner: owner,
		LeaseFor: cfg.OutboxLeaseDuration, BatchSize: int32(cfg.OutboxBatchSize),
	}
	go runRelay(ctx, relay, cfg.OutboxPollInterval)

	worker := email.Worker{Sender: sender, SendTimeout: cfg.EmailSendTimeout}
	var workers sync.WaitGroup
	for i := 0; i < cfg.EmailWorkerCount; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			consumer := rabbitmq.Consumer{Conn: conn, Queue: rabbitmq.SendQueue, Prefetch: cfg.EmailWorkerPrefetch}
			if err := consumer.Run(ctx, handleDelivery(worker, publisher)); err != nil && !errors.Is(err, context.Canceled) {
				logger.Error("email consumer stopped", zap.Error(err))
			}
		}()
	}
	logger.Info("email worker started", zap.Int("workers", cfg.EmailWorkerCount), zap.String("relay_owner", owner))
	<-ctx.Done()
	workers.Wait()
	logger.Info("email worker stopped")
}

func runRelay(ctx context.Context, relay outbox.Relay, interval time.Duration) {
	publish := func() {
		count, err := relay.RunOnce(ctx)
		if err != nil {
			logger.Error("outbox relay iteration failed", zap.Error(err))
			return
		}
		if count > 0 {
			logger.Info("outbox records relayed", zap.Int("count", count))
		}
	}
	publish()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			publish()
		}
	}
}

func handleDelivery(worker email.Worker, publisher rabbitmq.Publisher) func(context.Context, rabbitmq.Delivery) rabbitmq.Action {
	return func(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action {
		job, err := email.Unmarshal(delivery.Body)
		if err != nil {
			return publishFailure(ctx, publisher, delivery, "malformed_job: "+err.Error())
		}
		if err := worker.Deliver(ctx, job); err == nil {
			return rabbitmq.Ack
		} else {
			attempt := rabbitmq.RetryCount(delivery.Headers)
			if email.FailureOf(err) == email.FailurePermanent || attempt >= len(rabbitmq.RetryTiers) {
				return publishFailure(ctx, publisher, delivery, err.Error())
			}
			headers := rabbitmq.CopyHeaders(delivery.Headers)
			headers["x-retry-count"] = int32(attempt + 1)
			tier := rabbitmq.RetryTiers[attempt]
			if publishErr := publisher.PublishMessage(ctx, rabbitmq.Message{
				Exchange: rabbitmq.RetryExchange, RoutingKey: tier.RoutingKey, MessageID: job.ID,
				Body: delivery.Body, Headers: headers,
			}); publishErr != nil {
				logger.Error("schedule email retry", zap.Error(publishErr), zap.String("message_id", job.ID))
				return rabbitmq.Requeue
			}
			logger.Warn("email delivery failed; retry scheduled", zap.Error(err), zap.String("message_id", job.ID), zap.Int("attempt", attempt+1))
			return rabbitmq.Ack
		}
	}
}

func publishFailure(ctx context.Context, publisher rabbitmq.Publisher, delivery rabbitmq.Delivery, reason string) rabbitmq.Action {
	headers := rabbitmq.CopyHeaders(delivery.Headers)
	headers["x-failure-reason"] = reason
	if err := publisher.PublishMessage(ctx, rabbitmq.Message{
		Exchange: rabbitmq.DLXExchange, RoutingKey: rabbitmq.FailedRoutingKey, MessageID: delivery.MessageID,
		Body: delivery.Body, Headers: headers,
	}); err != nil {
		logger.Error("publish email failure to DLQ", zap.Error(err), zap.String("message_id", delivery.MessageID))
		return rabbitmq.Requeue
	}
	logger.Error("email sent to DLQ", zap.String("message_id", delivery.MessageID), zap.String("reason", reason))
	return rabbitmq.Ack
}

func relayOwner() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "email-worker"
	}
	return fmt.Sprintf("%s-%s", host, uuid.NewString())
}
