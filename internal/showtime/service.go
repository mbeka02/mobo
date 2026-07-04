package showtime

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("start time must be before end time")
)

// Service defines the business operations for the showtime domain.
type Service interface {
	CreateShowtime(ctx context.Context, req CreateShowtimeRequest) (*Showtime, error)
	GetShowtime(ctx context.Context, id int64) (*Showtime, error)
	ListShowtimesByMovie(ctx context.Context, movieId int64) ([]Showtime, error)
	ListShowtimesAdmin(ctx context.Context, limit, offset int32) ([]Showtime, error)
	UpdateShowtime(ctx context.Context, id int64, req UpdateShowtimeRequest) (*Showtime, error)
	DeleteShowtime(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

// NewService creates a new showtime service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateShowtime(ctx context.Context, req CreateShowtimeRequest) (*Showtime, error) {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	if !end.After(start) {
		return nil, ErrInvalidTimeRange
	}

	return s.repo.Create(ctx, req)
}

func (s *service) GetShowtime(ctx context.Context, id int64) (*Showtime, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListShowtimesByMovie(ctx context.Context, movieId int64) ([]Showtime, error) {
	return s.repo.ListByMovie(ctx, movieId)
}

func (s *service) ListShowtimesAdmin(ctx context.Context, limit, offset int32) ([]Showtime, error) {
	return s.repo.ListAdmin(ctx, limit, offset)
}

func (s *service) UpdateShowtime(ctx context.Context, id int64, req UpdateShowtimeRequest) (*Showtime, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *service) DeleteShowtime(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
