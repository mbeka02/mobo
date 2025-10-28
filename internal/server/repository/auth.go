package repository

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type CreateLocalUserParams struct {
	Email           string
	FullName        string
	PasswordHash    string
	TelephoneNumber string
}

type CreateOAuthUserParams struct {
	Email           string
	FullName        string
	AuthProvider    string
	ProviderUserId  string
	ProfileImageUrl string
}

type AuthRepository interface {
	CreateLocalUser(ctx context.Context, user CreateLocalUserParams) (*model.User, error)
	CreateOAuthUser(ctx context.Context, user CreateOAuthUserParams) (*model.User, error)
	GetUserByProvider(ctx context.Context, provider, providerUserID string) (*model.User, error)
}

type authRepository struct {
	store *database.Store
}

func NewAuthRepository(store *database.Store) AuthRepository {
	return &authRepository{store}
}

func (r *authRepository) GetUserByProvider(ctx context.Context, provider, providerUserID string) (*model.User, error) {
	row, err := r.store.GetUserByProvider(ctx, database.GetUserByProviderParams{
		AuthProvider:   &provider,
		ProviderUserID: &providerUserID,
	})
	if err != nil {
		return nil, err
	}

	return model.FromGetUserByProviderRow(&row), nil
}

func (r *authRepository) CreateOAuthUser(ctx context.Context, user CreateOAuthUserParams) (*model.User, error) {
	dbUser, err := r.store.CreateOAuthUser(ctx, database.CreateOAuthUserParams{
		Email:           user.Email,
		FullName:        user.FullName,
		AuthProvider:    &user.AuthProvider,
		ProviderUserID:  &user.ProviderUserId,
		ProfileImageUrl: &user.ProfileImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return model.FromDatabaseUser(&dbUser), nil
}

func (r *authRepository) CreateLocalUser(ctx context.Context, user CreateLocalUserParams) (*model.User, error) {
	dbUser, err := r.store.CreateLocalUser(ctx, database.CreateLocalUserParams{
		Email:           user.Email,
		FullName:        user.FullName,
		PasswordHash:    &user.PasswordHash,
		TelephoneNumber: user.TelephoneNumber,
	})
	if err != nil {
		return nil, err
	}

	return model.FromDatabaseUser(&dbUser), nil
}
