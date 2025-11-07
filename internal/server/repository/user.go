package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	store *database.Store
}

func NewUserRepository(store *database.Store) UserRepository {
	return &userRepository{store}
}

func (ur *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row, err := ur.store.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.FromGetUserByID(&row), nil
}

func (ur *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row, err := ur.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return model.FromGetUserByEmailRow(&row), nil
}
