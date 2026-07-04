package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

var (
	ErrEmailAlreadyExists = errors.New("a user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrOAuthOnlyAccount   = errors.New("this account uses a social login provider, please log in with your provider or set a password in account settings")
)

// Service defines the business operations for the user domain.
type Service interface {
	CreateOrLoginOAuthUser(ctx context.Context, data OAuthUserData) (*User, error)
	RegisterLocalUser(ctx context.Context, email, fullName, password, telephone string) (*User, error)
	LoginLocalUser(ctx context.Context, email, password string) (*User, error)
}

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateOrLoginOAuthUser(ctx context.Context, data OAuthUserData) (*User, error) {
	// 1. Check if identity already exists
	existingUser, err := s.repo.GetByProvider(ctx, data.Provider, data.ProviderUserID)
	if err == nil {
		logger.DebugCtx(ctx, "existing user found, logging in",
			zap.String("user_id", existingUser.ID.String()),
			zap.String("provider", data.Provider),
		)
		return existingUser, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("error checking existing identity: %w", err)
	}

	// 2. Identity not found, check if user with same email exists
	userByEmail, err := s.repo.GetByEmail(ctx, data.Email)
	if err == nil {
		// 3. User exists, link new identity to this user
		logger.DebugCtx(ctx, "user with same email found, linking new identity",
			zap.String("user_id", userByEmail.ID.String()),
			zap.String("provider", data.Provider),
			zap.String("email", data.Email),
		)

		err = s.repo.LinkIdentity(ctx, userByEmail.ID, data.Provider, data.ProviderUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to link identity: %w", err)
		}

		return userByEmail, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("error checking user by email: %w", err)
	}

	// 4. No user found, create new user AND identity in a single transaction
	logger.DebugCtx(ctx, "no existing user found, creating new user and identity",
		zap.String("provider", data.Provider),
		zap.String("email", data.Email),
	)

	fullName := fmt.Sprintf("%s %s", data.FirstName, data.LastName)
	newUser, err := s.repo.CreateWithIdentity(ctx, CreateUserParams{
		Email:           data.Email,
		FullName:        fullName,
		ProfileImageUrl: data.AvatarURL,
	}, data.Provider, data.ProviderUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with identity: %w", err)
	}

	return newUser, nil
}

func (s *service) RegisterLocalUser(ctx context.Context, email, fullName, password, telephone string) (*User, error) {
	// Check if email is already taken
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		logger.WarnCtx(ctx, "registration attempt with existing email",
			zap.String("email", email),
		)
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create the user and link local identity in a single transaction
	user, err := s.repo.CreateLocalWithIdentity(ctx, email, fullName, hashedPassword, telephone)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create local user",
			zap.Error(err),
			zap.String("email", email),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.InfoCtx(ctx, "local user registered successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", email),
	)

	return user, nil
}

func (s *service) LoginLocalUser(ctx context.Context, email, password string) (*User, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			logger.WarnCtx(ctx, "login attempt with non-existent email",
				zap.String("email", email),
			)
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	// Ensure user has a password (not an OAuth-only account)
	if user.PasswordHash == nil || *user.PasswordHash == "" {
		logger.WarnCtx(ctx, "login attempt on OAuth-only account",
			zap.String("email", email),
		)
		return nil, ErrOAuthOnlyAccount
	}

	// Compare password
	if err := auth.ComparePassword(password, *user.PasswordHash); err != nil {
		logger.WarnCtx(ctx, "login attempt with wrong password",
			zap.String("email", email),
		)
		return nil, ErrInvalidCredentials
	}

	logger.InfoCtx(ctx, "local user logged in successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", email),
	)

	return user, nil
}
