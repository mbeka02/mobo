package user

import (
	"time"

	"github.com/google/uuid"
)

// Role constants for the user domain.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents a user in the system.
type User struct {
	ID              uuid.UUID
	Email           string
	Role            string
	TelephoneNumber *string
	PasswordHash    *string
	FullName        string
	ProfileImageURL *string
	UserName        *string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	VerifiedAt      *time.Time
}

// ToResponse converts a User to a UserResponse.
func (u *User) ToResponse() UserResponse {
	var verifiedAt time.Time
	if u.VerifiedAt != nil {
		verifiedAt = *u.VerifiedAt
	}

	var updatedAt time.Time
	if u.UpdatedAt != nil {
		updatedAt = *u.UpdatedAt
	}

	return UserResponse{
		UserId:          u.ID.String(),
		Fullname:        u.FullName,
		Email:           u.Email,
		TelephoneNumber: u.TelephoneNumber,
		ProfileImageURL: u.ProfileImageURL,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       updatedAt,
		VerifiedAt:      verifiedAt,
	}
}

// Request/Response models

// CreateLocalUserRequest represents the request to register a local user.
type CreateLocalUserRequest struct {
	Fullname        string `json:"full_name" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	TelephoneNumber string `json:"telephone_number" validate:"required,max=15"`
	Password        string `json:"password" validate:"required,min=8"`
}

// UserResponse represents the API response for a user.
type UserResponse struct {
	UserId          string    `json:"user_id"`
	Fullname        string    `json:"full_name"`
	Email           string    `json:"email"`
	TelephoneNumber *string   `json:"telephone_number"`
	ProfileImageURL *string   `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	VerifiedAt      time.Time `json:"verified_at"`
}

// LoginRequest represents the request to log in.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// OAuthUserData contains the user data received from an OAuth provider.
type OAuthUserData struct {
	Email          string
	FirstName      string
	LastName       string
	Provider       string
	ProviderUserID string
	AvatarURL      string
}
