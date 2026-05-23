package repository

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type MovieRepository interface {
	AddMovie(ctx context.Context, params database.AddMovieParams) (*model.Movie, error)
	GetMovieById(ctx context.Context, id int64) (*model.Movie, error)
	GetMoviesAdmin(ctx context.Context, limit, offset int32) ([]model.Movie, error)
	GetMoviesPublic(ctx context.Context, limit, offset int32) ([]model.Movie, error)
	UpdateMovie(ctx context.Context, params database.UpdateMovieParams) (*model.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type movieRepository struct {
	store *database.Store
}

func NewMovieRepository(store *database.Store) MovieRepository {
	return &movieRepository{store}
}

func (r *movieRepository) AddMovie(ctx context.Context, params database.AddMovieParams) (*model.Movie, error) {
	movie, err := r.store.AddMovie(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseMovie(&movie), nil
}

func (r *movieRepository) GetMovieById(ctx context.Context, id int64) (*model.Movie, error) {
	movie, err := r.store.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseGetMovieByIdRow(&movie), nil
}

func (r *movieRepository) GetMoviesAdmin(ctx context.Context, limit, offset int32) ([]model.Movie, error) {
	movies, err := r.store.GetMoviesAdmin(ctx, database.GetMoviesAdminParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]model.Movie, 0, len(movies))
	for _, m := range movies {
		res = append(res, *model.FromDatabaseGetMoviesAdminRow(&m))
	}
	return res, nil
}

func (r *movieRepository) GetMoviesPublic(ctx context.Context, limit, offset int32) ([]model.Movie, error) {
	movies, err := r.store.GetMoviesPublic(ctx, database.GetMoviesPublicParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]model.Movie, 0, len(movies))
	for _, m := range movies {
		res = append(res, *model.FromDatabaseMovie(&m))
	}
	return res, nil
}

func (r *movieRepository) UpdateMovie(ctx context.Context, params database.UpdateMovieParams) (*model.Movie, error) {
	movie, err := r.store.UpdateMovie(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseMovie(&movie), nil
}

func (r *movieRepository) DeleteMovie(ctx context.Context, id int64) error {
	return r.store.DeleteMovie(ctx, id)
}
