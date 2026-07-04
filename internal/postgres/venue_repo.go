package postgres

import (
	"context"
	"time"

	"github.com/mbeka02/ticketing-service/internal/dbgen"
	"github.com/mbeka02/ticketing-service/internal/venue"
)

type venueRepo struct {
	store *Store
}

// NewVenueRepository creates a new postgres venue repository.
func NewVenueRepository(store *Store) venue.Repository {
	return &venueRepo{store}
}

func (r *venueRepo) Create(ctx context.Context, req venue.CreateVenueRequest) (*venue.Venue, error) {
	dbVenue, err := r.store.CreateVenue(ctx, dbgen.CreateVenueParams{
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		TotalSeats: req.TotalSeats,
	})
	if err != nil {
		return nil, err
	}
	return fromDatabaseVenue(&dbVenue), nil
}

func (r *venueRepo) GetByID(ctx context.Context, id int32) (*venue.Venue, error) {
	dbVenue, err := r.store.GetVenueById(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromDatabaseVenue(&dbVenue), nil
}

func (r *venueRepo) List(ctx context.Context) ([]venue.Venue, error) {
	venues, err := r.store.GetVenues(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]venue.Venue, 0, len(venues))
	for _, v := range venues {
		res = append(res, *fromDatabaseVenue(&v))
	}
	return res, nil
}

func fromDatabaseVenue(dbVenue *dbgen.Venue) *venue.Venue {
	var updatedAt *time.Time
	if dbVenue.UpdatedAt.Valid {
		updatedAt = &dbVenue.UpdatedAt.Time
	}

	return &venue.Venue{
		ID:         dbVenue.ID,
		Name:       dbVenue.Name,
		Address:    dbVenue.Address,
		City:       dbVenue.City,
		TotalSeats: dbVenue.TotalSeats,
		CreatedAt:  dbVenue.CreatedAt,
		UpdatedAt:  updatedAt,
	}
}
