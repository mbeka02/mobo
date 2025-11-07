package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/ticketing-service/internal/database"
)

type User struct {
	ID              uuid.UUID
	Email           string
	TelephoneNumber *string
	FullName        string
	ProfileImageURL *string
	UserName        *string
	AuthProvider    *string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	VerifiedAt      *time.Time
}

func FromDatabaseUser(dbUser *database.User) *User {
	return &User{
		ID:              dbUser.ID,
		Email:           dbUser.Email,
		TelephoneNumber: dbUser.TelephoneNumber,
		FullName:        dbUser.FullName,
		ProfileImageURL: dbUser.ProfileImageUrl,
		UserName:        dbUser.UserName,
		AuthProvider:    dbUser.AuthProvider,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       &dbUser.UpdatedAt.Time,
		VerifiedAt:      &dbUser.VerifiedAt.Time,
	}
}

func FromGetUserByEmailRow(row *database.GetUserByEmailRow) *User {
	return &User{
		ID:              row.ID,
		Email:           row.Email,
		TelephoneNumber: row.TelephoneNumber,
		FullName:        row.FullName,
		ProfileImageURL: row.ProfileImageUrl,
		UserName:        row.UserName,
		AuthProvider:    row.AuthProvider,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       &row.UpdatedAt.Time,
		VerifiedAt:      &row.VerifiedAt.Time,
	}
}

func FromGetUserByID(row *database.GetUserByIdRow) *User {
	return &User{
		ID:              row.ID,
		Email:           row.Email,
		TelephoneNumber: row.TelephoneNumber,
		FullName:        row.FullName,
		ProfileImageURL: row.ProfileImageUrl,
		UserName:        row.UserName,
		AuthProvider:    row.AuthProvider,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       &row.UpdatedAt.Time,
		VerifiedAt:      &row.VerifiedAt.Time,
	}
}

func FromGetUserByProviderRow(row *database.GetUserByProviderRow) *User {
	return &User{
		ID:              row.ID,
		Email:           row.Email,
		TelephoneNumber: row.TelephoneNumber,
		FullName:        row.FullName,
		ProfileImageURL: row.ProfileImageUrl,
		UserName:        row.UserName,
		AuthProvider:    row.AuthProvider,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       &row.UpdatedAt.Time,
		VerifiedAt:      &row.VerifiedAt.Time,
	}
}

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
type CreateLocalUserRequest struct {
	Fullname        string `json:"full_name" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	TelephoneNumber string `json:"telephone_number" validate:"required,max=15"`
	Password        string `json:"password" validate:"required,min=8"`
}

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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
