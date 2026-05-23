package service

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/model"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
)

type VenueService interface {
	CreateVenue(ctx context.Context, req model.CreateVenueRequest) (*model.Venue, error)
	GetVenue(ctx context.Context, id int32) (*model.Venue, error)
	ListVenues(ctx context.Context) ([]model.Venue, error)
}

type venueService struct {
	repo repository.VenueRepository
}

func NewVenueService(repo repository.VenueRepository) VenueService {
	return &venueService{repo}
}

func (s *venueService) CreateVenue(ctx context.Context, req model.CreateVenueRequest) (*model.Venue, error) {
	return s.repo.CreateVenue(ctx, req.ToParams())
}

func (s *venueService) GetVenue(ctx context.Context, id int32) (*model.Venue, error) {
	return s.repo.GetVenueById(ctx, id)
}

func (s *venueService) ListVenues(ctx context.Context) ([]model.Venue, error) {
	return s.repo.GetVenues(ctx)
}
