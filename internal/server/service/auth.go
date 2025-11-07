package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/mbeka02/ticketing-service/internal/model"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
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
		logger.InfoCtx(ctx, "existing OAuth user found, logging in",
			zap.String("user_id", existingUser.ID.String()),
			zap.String("provider", data.Provider),
		)
		return existingUser, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		logger.InfoCtx(ctx, "no existing user found, creating new OAuth user",
			zap.String("provider", data.Provider),
			zap.String("email", data.Email),
		)

		fullName := fmt.Sprintf("%s %s", data.FirstName, data.LastName)
		newUser, createErr := s.repo.CreateOAuthUser(ctx, repository.CreateOAuthUserParams{
			Email:           data.Email,
			FullName:        fullName,
			AuthProvider:    data.Provider,
			ProviderUserId:  data.ProviderUserID,
			ProfileImageUrl: data.AvatarURL,
		})
		if createErr != nil {
			logger.ErrorCtx(ctx, "failed to create OAuth user",
				zap.Error(createErr),
				zap.String("provider", data.Provider),
				zap.String("email", data.Email),
			)
			return nil, fmt.Errorf("error creating oauth user: %w", createErr)
		}

		logger.InfoCtx(ctx, "successfully created new OAuth user",
			zap.String("user_id", newUser.ID.String()),
			zap.String("provider", data.Provider),
		)
		return newUser, nil
	}

	logger.ErrorCtx(ctx, "unexpected error checking existing user",
		zap.Error(err),
		zap.String("provider", data.Provider),
	)

	return nil, fmt.Errorf("error checking existing user: %w", err)
}
