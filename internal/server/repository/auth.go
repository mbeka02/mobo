package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type CreateUserParams struct {
	Email           string
	FullName        string
	ProfileImageUrl string
	VerifiedAt      time.Time
}

type AuthRepository interface {
	GetUserByProvider(ctx context.Context, provider, providerUserID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user CreateUserParams) (*model.User, error)
	LinkIdentityToUser(ctx context.Context, userID uuid.UUID, provider, providerUserID string) error
	CreateUserWithIdentity(ctx context.Context, user CreateUserParams, provider, providerUserID string) (*model.User, error)
	CreateLocalUser(ctx context.Context, email, fullName, passwordHash, telephone string) (*model.User, error)
}

type authRepository struct {
	store *database.Store
}

func NewAuthRepository(store *database.Store) AuthRepository {
	return &authRepository{store}
}

func (r *authRepository) GetUserByProvider(ctx context.Context, provider, providerUserID string) (*model.User, error) {
	dbUser, err := r.store.GetUserByProvider(ctx, database.GetUserByProviderParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if err != nil {
		return nil, err
	}

	return model.FromGetUserByProviderRow(&dbUser), nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return model.FromGetUserByEmailRow(&dbUser), nil
}

func (r *authRepository) CreateUser(ctx context.Context, user CreateUserParams) (*model.User, error) {
	dbUser, err := r.store.CreateUser(ctx, database.CreateUserParams{
		Email:           user.Email,
		FullName:        user.FullName,
		ProfileImageUrl: &user.ProfileImageUrl,
		VerifiedAt:      pgtype.Timestamptz{Time: user.VerifiedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseUser(&dbUser), nil
}

func (r *authRepository) LinkIdentityToUser(ctx context.Context, userID uuid.UUID, provider, providerUserID string) error {
	_, err := r.store.LinkIdentityToUser(ctx, database.LinkIdentityToUserParams{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	return err
}

func (r *authRepository) CreateUserWithIdentity(ctx context.Context, user CreateUserParams, provider, providerUserID string) (*model.User, error) {
	var createdUser *model.User
	err := r.store.ExecTx(ctx, func(q *database.Queries) error {
		dbUser, err := q.CreateUser(ctx, database.CreateUserParams{
			Email:           user.Email,
			FullName:        user.FullName,
			ProfileImageUrl: &user.ProfileImageUrl,
			VerifiedAt:      pgtype.Timestamptz{Time: user.VerifiedAt, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create user in transaction: %w", err)
		}

		_, err = q.LinkIdentityToUser(ctx, database.LinkIdentityToUserParams{
			UserID:         dbUser.ID,
			Provider:       provider,
			ProviderUserID: providerUserID,
		})
		if err != nil {
			return fmt.Errorf("failed to link identity in transaction: %w", err)
		}

		createdUser = model.FromDatabaseUser(&dbUser)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (r *authRepository) CreateLocalUser(ctx context.Context, email, fullName, passwordHash, telephone string) (*model.User, error) {
	dbUser, err := r.store.CreateLocalUser(ctx, database.CreateLocalUserParams{
		Email:           email,
		FullName:        fullName,
		PasswordHash:    &passwordHash,
		TelephoneNumber: &telephone,
	})
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseUser(&dbUser), nil
}
