package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	RegisterLocalUser(ctx context.Context, email, fullName, password, telephone string) (*model.User, error)
	LoginLocalUser(ctx context.Context, email, password string) (*model.User, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) CreateOrLoginOAuthUser(ctx context.Context, data OAuthUserData) (*model.User, error) {
	// 1. Check if identity already exists
	existingUser, err := s.repo.GetUserByProvider(ctx, data.Provider, data.ProviderUserID)
	if err == nil {
		logger.DebugCtx(ctx, "existing user found, logging in",
			zap.String("user_id", existingUser.ID.String()),
			zap.String("provider", data.Provider),
		)
		return existingUser, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error checking existing identity: %w", err)
	}

	// 2. Identity not found, check if user with same email exists
	userByEmail, err := s.repo.GetUserByEmail(ctx, data.Email)
	if err == nil {
		// 3. User exists, link new identity to this user
		logger.DebugCtx(ctx, "user with same email found, linking new identity",
			zap.String("user_id", userByEmail.ID.String()),
			zap.String("provider", data.Provider),
			zap.String("email", data.Email),
		)

		err = s.repo.LinkIdentityToUser(ctx, userByEmail.ID, data.Provider, data.ProviderUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to link identity: %w", err)
		}

		return userByEmail, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error checking user by email: %w", err)
	}

	// 4. No user found, create new user AND identity in a single transaction
	logger.DebugCtx(ctx, "no existing user found, creating new user and identity",
		zap.String("provider", data.Provider),
		zap.String("email", data.Email),
	)

	fullName := fmt.Sprintf("%s %s", data.FirstName, data.LastName)
	newUser, err := s.repo.CreateUserWithIdentity(ctx, repository.CreateUserParams{
		Email:           data.Email,
		FullName:        fullName,
		ProfileImageUrl: data.AvatarURL,
		VerifiedAt:      time.Now(),
	}, data.Provider, data.ProviderUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with identity: %w", err)
	}

	return newUser, nil
}

func (s *authService) RegisterLocalUser(ctx context.Context, email, fullName, password, telephone string) (*model.User, error) {
	// TODO: Implement password hashing
	// TODO: Check if email already exists
	// TODO: Create user and identity
	return nil, errors.New("not implemented")
}

func (s *authService) LoginLocalUser(ctx context.Context, email, password string) (*model.User, error) {
	// TODO: Get user by email
	// TODO: Verify password hash
	// TODO: Check for 'local' identity
	return nil, errors.New("not implemented")
}
