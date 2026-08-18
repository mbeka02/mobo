package resend

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbeka02/ticketing-service/internal/email"
	resendsdk "github.com/resend/resend-go/v3"
)

type Client interface {
	SendWithOptions(context.Context, *resendsdk.SendEmailRequest, *resendsdk.SendEmailOptions) (*resendsdk.SendEmailResponse, error)
}

type Sender struct {
	client Client
	from   string
}

func New(apiKey, from string) (*Sender, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("RESEND_API_KEY is required")
	}
	if from == "" {
		return nil, fmt.Errorf("EMAIL_FROM is required")
	}
	return &Sender{client: resendsdk.NewClient(apiKey).Emails, from: from}, nil
}

func NewWithClient(client Client, from string) *Sender { return &Sender{client: client, from: from} }

func (s *Sender) Send(ctx context.Context, message email.Message) error {
	if s.client == nil {
		return email.Permanent(errors.New("resend client is required"))
	}
	_, err := s.client.SendWithOptions(ctx, &resendsdk.SendEmailRequest{
		From: s.from, To: []string{message.To}, Subject: message.Subject, Html: message.HTML, Text: message.Text,
	}, &resendsdk.SendEmailOptions{IdempotencyKey: message.ID})
	if err == nil {
		return nil
	}
	if errors.Is(err, resendsdk.ErrRateLimit) {
		return email.RateLimited(err)
	}
	var missing *resendsdk.MissingRequiredFieldsError
	if errors.As(err, &missing) {
		return email.Permanent(err)
	}
	return email.Transient(err)
}
