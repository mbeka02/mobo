package movie

import "context"

// Repository defines the data access contract for the movie domain.
type Repository interface {
	Add(ctx context.Context, req AddMovieRequest) (*Movie, error)
	GetByID(ctx context.Context, id int64) (*Movie, error)
	ListAdmin(ctx context.Context, limit, offset int32) ([]Movie, error)
	ListPublic(ctx context.Context, limit, offset int32) ([]Movie, error)
	Update(ctx context.Context, id int64, req UpdateMovieRequest) (*Movie, error)
	Delete(ctx context.Context, id int64) error
}
