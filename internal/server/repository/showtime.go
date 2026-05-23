package repository

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type ShowtimeRepository interface {
	CreateShowtime(ctx context.Context, params database.CreateShowtimeParams) (*model.Showtime, error)
	GetShowtimeById(ctx context.Context, id int64) (*model.Showtime, error)
	GetShowtimesByMovie(ctx context.Context, movieId int64) ([]model.Showtime, error)
	GetShowtimesAdmin(ctx context.Context, limit, offset int32) ([]model.Showtime, error)
	UpdateShowtime(ctx context.Context, params database.UpdateShowtimeParams) (*model.Showtime, error)
	DeleteShowtime(ctx context.Context, id int64) error
}

type showtimeRepository struct {
	store *database.Store
}

func NewShowtimeRepository(store *database.Store) ShowtimeRepository {
	return &showtimeRepository{store}
}

func (r *showtimeRepository) CreateShowtime(ctx context.Context, params database.CreateShowtimeParams) (*model.Showtime, error) {
	showtime, err := r.store.CreateShowtime(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseShowtime(&showtime), nil
}

func (r *showtimeRepository) GetShowtimeById(ctx context.Context, id int64) (*model.Showtime, error) {
	showtime, err := r.store.GetShowtimeById(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseShowtime(&showtime), nil
}

func (r *showtimeRepository) GetShowtimesByMovie(ctx context.Context, movieId int64) ([]model.Showtime, error) {
	showtimes, err := r.store.GetShowtimesByMovie(ctx, movieId)
	if err != nil {
		return nil, err
	}

	res := make([]model.Showtime, 0, len(showtimes))
	for _, s := range showtimes {
		res = append(res, *model.FromDatabaseGetShowtimesByMovieRow(&s))
	}
	return res, nil
}

func (r *showtimeRepository) GetShowtimesAdmin(ctx context.Context, limit, offset int32) ([]model.Showtime, error) {
	showtimes, err := r.store.GetShowtimesAdmin(ctx, database.GetShowtimesAdminParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]model.Showtime, 0, len(showtimes))
	for _, s := range showtimes {
		res = append(res, *model.FromDatabaseGetShowtimesAdminRow(&s))
	}
	return res, nil
}

func (r *showtimeRepository) UpdateShowtime(ctx context.Context, params database.UpdateShowtimeParams) (*model.Showtime, error) {
	showtime, err := r.store.UpdateShowtime(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseShowtime(&showtime), nil
}

func (r *showtimeRepository) DeleteShowtime(ctx context.Context, id int64) error {
	return r.store.DeleteShowtime(ctx, id)
}
