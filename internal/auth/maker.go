package auth

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	Create(userId uuid.UUID, email string, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
