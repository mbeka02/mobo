package venue

import "time"

// Venue represents a venue in the system.
type Venue struct {
	ID         int32
	Name       string
	Address    string
	City       string
	TotalSeats int32
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

// ToResponse converts a Venue to a VenueResponse.
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

// VenueResponse represents the API response for a venue.
type VenueResponse struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	TotalSeats int32     `json:"total_seats"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// CreateVenueRequest represents the request to create a venue.
type CreateVenueRequest struct {
	Name       string `json:"name" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	TotalSeats int32  `json:"total_seats" validate:"required,min=1"`
}
