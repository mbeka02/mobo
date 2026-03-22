package auth

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	Create(userId uuid.UUID, email string, tokenType TokenType, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
