package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mbeka02/ticketing-service/internal/model"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
)

type OAuthUserData struct {
	Email          string
	FirstName      string
	LastName       string
	Provider       string
	ProviderUserID string
	AvatarURL      string
}
type AuthService interface {
	CreateOrLoginOAuthUser(ctx context.Context, data OAuthUserData) (*model.User, error)
	CreateLocalUser()
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) CreateLocalUser() {
}

func (s *authService) CreateOrLoginOAuthUser(ctx context.Context, data OAuthUserData) (*model.User, error) {
	existingUser, err := s.repo.GetUserByProvider(ctx, data.Provider, data.ProviderUserID)

	if err == nil {
		return existingUser, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		fullName := fmt.Sprintf("%s %s", data.FirstName, data.LastName)

		newUser, createErr := s.repo.CreateOAuthUser(ctx, repository.CreateOAuthUserParams{
			Email:           data.Email,
			FullName:        fullName,
			AuthProvider:    data.Provider,
			ProviderUserId:  data.ProviderUserID,
			ProfileImageUrl: data.AvatarURL,
		})
		if createErr != nil {
			return nil, fmt.Errorf("error creating oauth user: (%w)", createErr)
		}

		return newUser, nil

	}
	// unexpected error when checking if user exists
	return nil, fmt.Errorf("error checking existing user: (%w)", err)
}
