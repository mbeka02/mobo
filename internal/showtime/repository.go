package showtime

import "context"

// Repository defines the data access contract for the showtime domain.
type Repository interface {
	Create(ctx context.Context, req CreateShowtimeRequest) (*Showtime, error)
	GetByID(ctx context.Context, id int64) (*Showtime, error)
	ListByMovie(ctx context.Context, movieId int64) ([]Showtime, error)
	ListAdmin(ctx context.Context, limit, offset int32) ([]Showtime, error)
	Update(ctx context.Context, id int64, req UpdateShowtimeRequest) (*Showtime, error)
	Delete(ctx context.Context, id int64) error
}
