package movie

import "context"

// Service defines the business operations for the movie domain.
type Service interface {
	AddMovie(ctx context.Context, req AddMovieRequest) (*Movie, error)
	GetMovie(ctx context.Context, id int64) (*Movie, error)
	ListMoviesAdmin(ctx context.Context, limit, offset int32) ([]Movie, error)
	ListMoviesPublic(ctx context.Context, limit, offset int32) ([]Movie, error)
	UpdateMovie(ctx context.Context, id int64, req UpdateMovieRequest) (*Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

// NewService creates a new movie service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddMovie(ctx context.Context, req AddMovieRequest) (*Movie, error) {
	return s.repo.Add(ctx, req)
}

func (s *service) GetMovie(ctx context.Context, id int64) (*Movie, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListMoviesAdmin(ctx context.Context, limit, offset int32) ([]Movie, error) {
	return s.repo.ListAdmin(ctx, limit, offset)
}

func (s *service) ListMoviesPublic(ctx context.Context, limit, offset int32) ([]Movie, error) {
	return s.repo.ListPublic(ctx, limit, offset)
}

func (s *service) UpdateMovie(ctx context.Context, id int64, req UpdateMovieRequest) (*Movie, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *service) DeleteMovie(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
