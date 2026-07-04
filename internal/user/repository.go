package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrNotFound is returned when a user is not found.
var ErrNotFound = errors.New("user not found")

// CreateUserParams contains the parameters for creating an OAuth user.
type CreateUserParams struct {
	Email           string
	FullName        string
	ProfileImageUrl string
	VerifiedAt      time.Time
}

// Repository defines the data access contract for the user domain.
type Repository interface {
	GetByProvider(ctx context.Context, provider, providerUserID string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	CreateWithIdentity(ctx context.Context, params CreateUserParams, provider, providerUserID string) (*User, error)
	CreateLocalWithIdentity(ctx context.Context, email, fullName, passwordHash, telephone string) (*User, error)
	LinkIdentity(ctx context.Context, userID uuid.UUID, provider, providerUserID string) error
}
