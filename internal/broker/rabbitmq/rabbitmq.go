package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/mbeka02/ticketing-service/internal/outbox"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// EventsExchange accepts newly relayed transactional-email jobs.
	EventsExchange = "email.events"
	// RetryExchange routes a failed job to one of the delayed retry queues.
	RetryExchange = "email.retry"
	// DLXExchange routes terminal failures to the queue operators inspect and replay.
	DLXExchange = "email.dlx"
	// SendQueue is the only queue consumed by the email worker.
	SendQueue = "email.send"
	// FailedQueue intentionally has no consumer; it is the email DLQ.
	FailedQueue      = "email.failed"
	SendRoutingKey   = "send"
	FailedRoutingKey = "failed"
)

// RetryTier describes a holding queue. On TTL expiry, RabbitMQ dead-letters
// the message back to SendQueue; no worker consumes retry queues directly.
type RetryTier struct {
	RoutingKey string
	Queue      string
	Delay      time.Duration
}

// RetryTiers is the fixed, escalating retry policy. The retry count header is
// the zero-based index into this slice, so changing the order changes delivery
// behavior for in-flight messages.
var RetryTiers = []RetryTier{
	{RoutingKey: "retry.1m", Queue: "email.retry.1m", Delay: time.Minute},
	{RoutingKey: "retry.5m", Queue: "email.retry.5m", Delay: 5 * time.Minute},
	{RoutingKey: "retry.15m", Queue: "email.retry.15m", Delay: 15 * time.Minute},
	{RoutingKey: "retry.1h", Queue: "email.retry.1h", Delay: time.Hour},
	{RoutingKey: "retry.6h", Queue: "email.retry.6h", Delay: 6 * time.Hour},
}

// Topology declares the exchanges, queues, and bindings used by the email
// pipeline. Queue/exchange declaration is idempotent as long as the deployed
// arguments remain identical, so every worker may safely call Setup at startup.
//
// QueueType defaults to quorum for production data safety. A local, single-node
// developer broker can opt into "classic" if it does not support quorum queues.
type Topology struct{ QueueType string }

// Setup creates durable broker resources. Durability preserves declarations
// across broker restarts; publishers must also mark individual messages as
// persistent, which Publisher.PublishMessage does below.
func (t Topology) Setup(conn *amqp.Connection) error {
	queueType := t.QueueType
	if queueType == "" {
		queueType = "quorum"
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open topology channel: %w", err)
	}
	defer ch.Close()
	for _, exchange := range []string{EventsExchange, RetryExchange, DLXExchange} {
		if err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
			return fmt.Errorf("declare exchange %s: %w", exchange, err)
		}
	}
	queueArgs := amqp.Table{"x-queue-type": queueType}
	if _, err := ch.QueueDeclare(SendQueue, true, false, false, false, queueArgs); err != nil {
		return fmt.Errorf("declare send queue: %w", err)
	}
	if err := ch.QueueBind(SendQueue, SendRoutingKey, EventsExchange, false, nil); err != nil {
		return fmt.Errorf("bind send queue: %w", err)
	}
	for _, tier := range RetryTiers {
		// TTL expiry is the delay mechanism. Dead-lettering returns the original
		// message to the main exchange using the same routing key it had initially.
		args := amqp.Table{
			"x-queue-type":              queueType,
			"x-message-ttl":             int32(tier.Delay.Milliseconds()),
			"x-dead-letter-exchange":    EventsExchange,
			"x-dead-letter-routing-key": SendRoutingKey,
		}
		if _, err := ch.QueueDeclare(tier.Queue, true, false, false, false, args); err != nil {
			return fmt.Errorf("declare retry queue %s: %w", tier.Queue, err)
		}
		if err := ch.QueueBind(tier.Queue, tier.RoutingKey, RetryExchange, false, nil); err != nil {
			return fmt.Errorf("bind retry queue %s: %w", tier.Queue, err)
		}
	}
	if _, err := ch.QueueDeclare(FailedQueue, true, false, false, false, queueArgs); err != nil {
		return fmt.Errorf("declare failed queue: %w", err)
	}
	if err := ch.QueueBind(FailedQueue, FailedRoutingKey, DLXExchange, false, nil); err != nil {
		return fmt.Errorf("bind failed queue: %w", err)
	}
	return nil
}

// Message is the transport-level representation used for normal publishes,
// retry scheduling, and DLQ routing. It deliberately carries no provider SDK
// types so the broker adapter remains independent of mailer implementations.
type Message struct {
	Exchange   string
	RoutingKey string
	MessageID  string
	Body       []byte
	Headers    map[string]any
}

// Publisher opens a short-lived channel per operation. AMQP channels are not
// safe for concurrent use, while a shared connection is safe for this pattern.
// A publish only succeeds after RabbitMQ sends a publisher confirmation.
type Publisher struct {
	Conn           *amqp.Connection
	ConfirmTimeout time.Duration
}

// Publish implements outbox.Publisher. The outbox record ID becomes the AMQP
// message ID and later the email provider idempotency key.
func (p Publisher) Publish(ctx context.Context, record outbox.Record) error {
	return p.PublishMessage(ctx, Message{
		Exchange: EventsExchange, RoutingKey: SendRoutingKey, MessageID: record.ID.String(), Body: record.Payload,
		Headers: map[string]any{"event_type": record.EventType, "schema_version": int32(record.SchemaVersion)},
	})
}

// PublishMessage persists one message and waits for a broker confirmation.
// Returning an error leaves the outbox record pending or tells the consumer to
// retain its original delivery, avoiding an acknowledge-before-retry-publish
// loss window.
func (p Publisher) PublishMessage(ctx context.Context, message Message) error {
	if p.Conn == nil {
		return fmt.Errorf("rabbitmq connection is required")
	}
	ch, err := p.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open publisher channel: %w", err)
	}
	defer ch.Close()
	// Confirm mode lets the caller distinguish a network write from a broker-
	// accepted publish. The latter is required before outbox completion or ack.
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("enable publisher confirms: %w", err)
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := ch.PublishWithContext(ctx, message.Exchange, message.RoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json", DeliveryMode: amqp.Persistent, MessageId: message.MessageID,
		Timestamp: time.Now().UTC(), Headers: amqp.Table(message.Headers), Body: message.Body,
	}); err != nil {
		return fmt.Errorf("publish message: %w", err)
	}
	timeout := p.ConfirmTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case confirmation := <-confirms:
		if !confirmation.Ack {
			return fmt.Errorf("broker negatively acknowledged publish")
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return fmt.Errorf("publisher confirm timed out after %s", timeout)
	}
}

// Delivery contains the data the application worker needs to decide whether a
// broker delivery should be acknowledged or requeued.
type Delivery struct {
	Body      []byte
	MessageID string
	Headers   map[string]any
}

// Action is returned by a consumer handler. Ack removes the current broker
// delivery; Requeue leaves it available for redelivery after a local/broker
// failure prevented safe retry or DLQ publication.
type Action int

const (
	Ack Action = iota
	Requeue
)

// Consumer owns one channel and delivers messages serially to handler. Run one
// Consumer per worker goroutine to gain concurrency without sharing a channel.
type Consumer struct {
	Conn     *amqp.Connection
	Queue    string
	Prefetch int
}

// Run consumes with manual acknowledgements. Prefetch bounds the number of
// in-flight messages per worker, applying backpressure when provider sending is
// slow. On shutdown, Cancel stops new deliveries and any unacknowledged work is
// eligible for broker redelivery.
func (c Consumer) Run(ctx context.Context, handler func(context.Context, Delivery) Action) error {
	if c.Conn == nil {
		return fmt.Errorf("rabbitmq connection is required")
	}
	queue := c.Queue
	if queue == "" {
		queue = SendQueue
	}
	prefetch := c.Prefetch
	if prefetch <= 0 {
		prefetch = 10
	}
	ch, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open consumer channel: %w", err)
	}
	defer ch.Close()
	if err := ch.Qos(prefetch, 0, false); err != nil {
		return fmt.Errorf("set prefetch: %w", err)
	}
	consumerTag := fmt.Sprintf("email-worker-%d", time.Now().UnixNano())
	deliveries, err := ch.Consume(queue, consumerTag, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume %s: %w", queue, err)
	}
	for {
		select {
		case <-ctx.Done():
			_ = ch.Cancel(consumerTag, false)
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}
			// The handler first publishes a confirmed retry or DLQ copy when needed;
			// only then is this source delivery acknowledged.
			action := handler(ctx, Delivery{Body: delivery.Body, MessageID: delivery.MessageId, Headers: map[string]any(delivery.Headers)})
			if action == Requeue {
				if err := delivery.Nack(false, true); err != nil {
					return fmt.Errorf("requeue message: %w", err)
				}
			} else if err := delivery.Ack(false); err != nil {
				return fmt.Errorf("ack message: %w", err)
			}
		}
	}
}

// RetryCount reads the application-managed header. AMQP number decoding can
// vary by publisher/client, hence the supported integer representations.
func RetryCount(headers map[string]any) int {
	if headers == nil {
		return 0
	}
	switch value := headers["x-retry-count"].(type) {
	case int32:
		return int(value)
	case int64:
		return int(value)
	case int:
		return value
	default:
		return 0
	}
}

// CopyHeaders preserves correlation and retry metadata while safely adding a
// new retry count or terminal failure reason for a republished message.
func CopyHeaders(headers map[string]any) map[string]any {
	copy := make(map[string]any, len(headers)+1)
	for key, value := range headers {
		copy[key] = value
	}
	return copy
}
