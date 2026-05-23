package model

import (
	"time"

	"github.com/mbeka02/ticketing-service/internal/database"
)

type Venue struct {
	ID         int32
	Name       string
	Address    string
	City       string
	TotalSeats int32
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func FromDatabaseVenue(dbVenue *database.Venue) *Venue {
	var updatedAt *time.Time
	if dbVenue.UpdatedAt.Valid {
		updatedAt = &dbVenue.UpdatedAt.Time
	}

	return &Venue{
		ID:         dbVenue.ID,
		Name:       dbVenue.Name,
		Address:    dbVenue.Address,
		City:       dbVenue.City,
		TotalSeats: dbVenue.TotalSeats,
		CreatedAt:  dbVenue.CreatedAt,
		UpdatedAt:  updatedAt,
	}
}

func (v *Venue) ToResponse() VenueResponse {
	var updatedAt time.Time
	if v.UpdatedAt != nil {
		updatedAt = *v.UpdatedAt
	}

	return VenueResponse{
		ID:         v.ID,
		Name:       v.Name,
		Address:    v.Address,
		City:       v.City,
		TotalSeats: v.TotalSeats,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  updatedAt,
	}
}

type VenueResponse struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	TotalSeats int32     `json:"total_seats"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type CreateVenueRequest struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	TotalSeats int32  `json:"total_seats" validate:"required,min=1"`
}

func (req *CreateVenueRequest) ToParams() database.CreateVenueParams {
	return database.CreateVenueParams{
		Name:       req.Name,
		Address:    req.Address,
		City:       req.City,
		TotalSeats: req.TotalSeats,
	}
}
