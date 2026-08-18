package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	EventUserCreated = "user.created.v1"
	TemplateWelcome  = "welcome.v1"
)

// Job is the versioned, provider-neutral payload transported through the queue.
type Job struct {
	ID            string            `json:"id"`
	EventType     string            `json:"event_type"`
	SchemaVersion int               `json:"schema_version"`
	AggregateType string            `json:"aggregate_type"`
	AggregateID   string            `json:"aggregate_id"`
	To            string            `json:"to"`
	Template      string            `json:"template"`
	Data          map[string]string `json:"data"`
	CreatedAt     time.Time         `json:"created_at"`
}

// Message is the fully rendered content passed to a provider adapter.
type Message struct {
	ID      string
	To      string
	Subject string
	HTML    string
	Text    string
}

// Sender is intentionally small so providers can be swapped without leaking SDK types.
type Sender interface {
	Send(context.Context, Message) error
}

type FailureKind string

const (
	FailurePermanent   FailureKind = "permanent"
	FailureTransient   FailureKind = "transient"
	FailureRateLimited FailureKind = "rate_limited"
)

type DeliveryError struct {
	Kind FailureKind
	Err  error
}

func (e *DeliveryError) Error() string { return e.Err.Error() }
func (e *DeliveryError) Unwrap() error { return e.Err }

func Permanent(err error) error   { return &DeliveryError{Kind: FailurePermanent, Err: err} }
func Transient(err error) error   { return &DeliveryError{Kind: FailureTransient, Err: err} }
func RateLimited(err error) error { return &DeliveryError{Kind: FailureRateLimited, Err: err} }

func FailureOf(err error) FailureKind {
	var deliveryErr *DeliveryError
	if errors.As(err, &deliveryErr) {
		return deliveryErr.Kind
	}
	return FailureTransient
}

func NewWelcomeJob(userID uuid.UUID, recipient, fullName string) Job {
	return Job{
		ID:            uuid.NewString(),
		EventType:     EventUserCreated,
		SchemaVersion: 1,
		AggregateType: "user",
		AggregateID:   userID.String(),
		To:            recipient,
		Template:      TemplateWelcome,
		Data:          map[string]string{"full_name": fullName},
		CreatedAt:     time.Now().UTC(),
	}
}

func Marshal(job Job) ([]byte, error) { return json.Marshal(job) }

func Unmarshal(payload []byte) (Job, error) {
	var job Job
	if err := json.Unmarshal(payload, &job); err != nil {
		return Job{}, fmt.Errorf("decode email job: %w", err)
	}
	if job.ID == "" || job.To == "" || job.Template == "" || job.EventType == "" {
		return Job{}, fmt.Errorf("invalid email job: id, event_type, to, and template are required")
	}
	return job, nil
}

// TODO: Work on better email templates
func Render(job Job) (Message, error) {
	switch job.Template {
	case TemplateWelcome:
		name := strings.TrimSpace(job.Data["full_name"])
		if name == "" {
			name = "there"
		}
		return Message{
			ID:      job.ID,
			To:      job.To,
			Subject: "Welcome to Mobo Ticketing",
			HTML:    fmt.Sprintf("<h1>Welcome to Mobo Ticketing, %s!</h1><p>Your account is ready.</p>", name),
			Text:    fmt.Sprintf("Welcome to Mobo Ticketing, %s! Your account is ready.", name),
		}, nil
	default:
		return Message{}, fmt.Errorf("unknown email template %q", job.Template)
	}
}

type Worker struct {
	Sender      Sender
	SendTimeout time.Duration
}

func (w Worker) Deliver(ctx context.Context, job Job) error {
	if w.Sender == nil {
		return Permanent(errors.New("email sender is required"))
	}
	message, err := Render(job)
	if err != nil {
		return Permanent(err)
	}
	timeout := w.SendTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	sendCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return w.Sender.Send(sendCtx, message)
}
