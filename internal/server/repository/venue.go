package repository

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
)

type VenueRepository interface {
	CreateVenue(ctx context.Context, params database.CreateVenueParams) (*model.Venue, error)
	GetVenueById(ctx context.Context, id int32) (*model.Venue, error)
	GetVenues(ctx context.Context) ([]model.Venue, error)
}

type venueRepository struct {
	store *database.Store
}

func NewVenueRepository(store *database.Store) VenueRepository {
	return &venueRepository{store}
}

func (r *venueRepository) CreateVenue(ctx context.Context, params database.CreateVenueParams) (*model.Venue, error) {
	venue, err := r.store.CreateVenue(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseVenue(&venue), nil
}

func (r *venueRepository) GetVenueById(ctx context.Context, id int32) (*model.Venue, error) {
	venue, err := r.store.GetVenueById(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.FromDatabaseVenue(&venue), nil
}

func (r *venueRepository) GetVenues(ctx context.Context) ([]model.Venue, error) {
	venues, err := r.store.GetVenues(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]model.Venue, 0, len(venues))
	for _, v := range venues {
		res = append(res, *model.FromDatabaseVenue(&v))
	}
	return res, nil
}
