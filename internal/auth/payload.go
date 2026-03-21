package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    uuid.UUID `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func NewPayload(userId uuid.UUID, email string, tokenType TokenType, duration time.Duration) *Payload {
	return &Payload{
		UserID:    userId,
		Email:     email,
		TokenType: tokenType,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}
