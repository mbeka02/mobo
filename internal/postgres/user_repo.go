package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/dbgen"
	"github.com/mbeka02/ticketing-service/internal/user"
)

type userRepo struct {
	store *Store
}

// NewUserRepository creates a new postgres user repository.
func NewUserRepository(store *Store) user.Repository {
	return &userRepo{store}
}

func (r *userRepo) GetByProvider(ctx context.Context, provider, providerUserID string) (*user.User, error) {
	dbUser, err := r.store.GetUserByProvider(ctx, dbgen.GetUserByProviderParams{
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return fromDatabaseUser(&dbUser), nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return fromDatabaseUser(&dbUser), nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	dbUser, err := r.store.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}
	return fromDatabaseUser(&dbUser), nil
}

func (r *userRepo) CreateWithIdentity(ctx context.Context, params user.CreateUserParams, provider, providerUserID string) (*user.User, error) {
	var createdUser *user.User
	err := r.store.ExecTx(ctx, func(q *dbgen.Queries) error {
		dbUser, err := q.CreateUser(ctx, dbgen.CreateUserParams{
			Email:           params.Email,
			FullName:        params.FullName,
			ProfileImageUrl: &params.ProfileImageUrl,
			VerifiedAt:      pgtype.Timestamptz{Time: params.VerifiedAt, Valid: !params.VerifiedAt.IsZero()},
		})
		if err != nil {
			return fmt.Errorf("failed to create user in transaction: %w", err)
		}

		_, err = q.LinkIdentityToUser(ctx, dbgen.LinkIdentityToUserParams{
			UserID:         dbUser.ID,
			Provider:       provider,
			ProviderUserID: providerUserID,
		})
		if err != nil {
			return fmt.Errorf("failed to link identity in transaction: %w", err)
		}

		createdUser = fromDatabaseUser(&dbUser)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (r *userRepo) CreateLocalWithIdentity(ctx context.Context, email, fullName, passwordHash, telephone string) (*user.User, error) {
	var createdUser *user.User
	err := r.store.ExecTx(ctx, func(q *dbgen.Queries) error {
		dbUser, err := q.CreateLocalUser(ctx, dbgen.CreateLocalUserParams{
			Email:           email,
			FullName:        fullName,
			PasswordHash:    &passwordHash,
			TelephoneNumber: &telephone,
		})
		if err != nil {
			return fmt.Errorf("failed to create local user in transaction: %w", err)
		}

		_, err = q.LinkIdentityToUser(ctx, dbgen.LinkIdentityToUserParams{
			UserID:         dbUser.ID,
			Provider:       "local",
			ProviderUserID: dbUser.ID.String(),
		})
		if err != nil {
			return fmt.Errorf("failed to link local identity in transaction: %w", err)
		}

		createdUser = fromDatabaseUser(&dbUser)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (r *userRepo) LinkIdentity(ctx context.Context, userID uuid.UUID, provider, providerUserID string) error {
	_, err := r.store.LinkIdentityToUser(ctx, dbgen.LinkIdentityToUserParams{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
	})
	return err
}

// fromDatabaseUser converts a dbgen.User to a user.User domain type.
func fromDatabaseUser(dbUser *dbgen.User) *user.User {
	var updatedAt *time.Time
	if dbUser.UpdatedAt.Valid {
		updatedAt = &dbUser.UpdatedAt.Time
	}
	var verifiedAt *time.Time
	if dbUser.VerifiedAt.Valid {
		verifiedAt = &dbUser.VerifiedAt.Time
	}

	return &user.User{
		ID:              dbUser.ID,
		Email:           dbUser.Email,
		Role:            string(dbUser.Role),
		TelephoneNumber: dbUser.TelephoneNumber,
		PasswordHash:    dbUser.PasswordHash,
		FullName:        dbUser.FullName,
		ProfileImageURL: dbUser.ProfileImageUrl,
		UserName:        dbUser.UserName,
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       updatedAt,
		VerifiedAt:      verifiedAt,
	}
}
