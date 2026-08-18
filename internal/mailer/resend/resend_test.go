package resend

import (
	"context"
	"errors"
	"testing"

	"github.com/mbeka02/ticketing-service/internal/email"
	resendsdk "github.com/resend/resend-go/v3"
	"github.com/stretchr/testify/require"
)

type fakeClient struct {
	request *resendsdk.SendEmailRequest
	options *resendsdk.SendEmailOptions
	err     error
}

func (c *fakeClient) SendWithOptions(_ context.Context, request *resendsdk.SendEmailRequest, options *resendsdk.SendEmailOptions) (*resendsdk.SendEmailResponse, error) {
	c.request, c.options = request, options
	return &resendsdk.SendEmailResponse{}, c.err
}

func TestSenderPassesMessageAndIdempotencyKey(t *testing.T) {
	client := &fakeClient{}
	sender := NewWithClient(client, "Mobo <bookings@example.com>")
	err := sender.Send(context.Background(), email.Message{ID: "job-1", To: "user@example.com", Subject: "Hello", HTML: "<p>Hello</p>", Text: "Hello"})
	require.NoError(t, err)
	require.Equal(t, "job-1", client.options.IdempotencyKey)
	require.Equal(t, []string{"user@example.com"}, client.request.To)
}

func TestSenderClassifiesRateLimit(t *testing.T) {
	client := &fakeClient{err: &resendsdk.RateLimitError{}}
	sender := NewWithClient(client, "Mobo <bookings@example.com>")
	err := sender.Send(context.Background(), email.Message{ID: "job-1", To: "user@example.com"})
	require.True(t, errors.Is(err, resendsdk.ErrRateLimit))
	require.Equal(t, email.FailureRateLimited, email.FailureOf(err))
}
