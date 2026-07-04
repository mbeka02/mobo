package venue

import "context"

// Service defines the business operations for the venue domain.
type Service interface {
	CreateVenue(ctx context.Context, req CreateVenueRequest) (*Venue, error)
	GetVenue(ctx context.Context, id int32) (*Venue, error)
	ListVenues(ctx context.Context) ([]Venue, error)
}

type service struct {
	repo Repository
}

// NewService creates a new venue service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateVenue(ctx context.Context, req CreateVenueRequest) (*Venue, error) {
	return s.repo.Create(ctx, req)
}

func (s *service) GetVenue(ctx context.Context, id int32) (*Venue, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) ListVenues(ctx context.Context) ([]Venue, error) {
	return s.repo.List(ctx)
}
