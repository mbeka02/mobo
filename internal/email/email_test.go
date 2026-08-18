package email

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeSender struct {
	message Message
	err     error
}

func (s *fakeSender) Send(_ context.Context, msg Message) error { s.message = msg; return s.err }

func TestRenderWelcome(t *testing.T) {
	job := NewWelcomeJob(uuid.New(), "reader@example.com", "Ada")
	message, err := Render(job)
	require.NoError(t, err)
	require.Equal(t, "reader@example.com", message.To)
	require.Contains(t, message.Subject, "Welcome")
	require.Contains(t, message.HTML, "Ada")
}

func TestWorkerClassifiesUnknownSenderErrorAsTransient(t *testing.T) {
	job := NewWelcomeJob(uuid.New(), "reader@example.com", "Ada")
	worker := Worker{Sender: &fakeSender{err: errors.New("network unavailable")}}
	err := worker.Deliver(context.Background(), job)
	require.Equal(t, FailureTransient, FailureOf(err))
}
