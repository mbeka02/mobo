package showtime

import "time"

// Showtime represents a showtime in the system.
type Showtime struct {
	ID             int64
	MovieID        int64
	StartTime      time.Time
	EndTime        time.Time
	AvailableSeats int32
	PricePerSeat   float64
	VenueID        int32
	CreatedAt      time.Time
	UpdatedAt      *time.Time

	// Enriched fields (populated by joins)
	MovieTitle *string
	VenueName  *string
	VenueCity  *string
}

// ToResponse converts a Showtime to a ShowtimeResponse.
func (s *Showtime) ToResponse() ShowtimeResponse {
	var updatedAt time.Time
	if s.UpdatedAt != nil {
		updatedAt = *s.UpdatedAt
	}

	return ShowtimeResponse{
		ID:             s.ID,
		MovieID:        s.MovieID,
		StartTime:      s.StartTime,
		EndTime:        s.EndTime,
		AvailableSeats: s.AvailableSeats,
		PricePerSeat:   s.PricePerSeat,
		VenueID:        s.VenueID,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      updatedAt,
		MovieTitle:     s.MovieTitle,
		VenueName:      s.VenueName,
		VenueCity:      s.VenueCity,
	}
}

// ShowtimeResponse represents the API response for a showtime.
type ShowtimeResponse struct {
	ID             int64     `json:"id"`
	MovieID        int64     `json:"movie_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	AvailableSeats int32     `json:"available_seats"`
	PricePerSeat   float64   `json:"price_per_seat"`
	VenueID        int32     `json:"venue_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	MovieTitle     *string   `json:"movie_title,omitempty"`
	VenueName      *string   `json:"venue_name,omitempty"`
	VenueCity      *string   `json:"venue_city,omitempty"`
}

// CreateShowtimeRequest represents the request to create a showtime.
type CreateShowtimeRequest struct {
	MovieID        int64   `json:"movie_id" validate:"required"`
	StartTime      string  `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime        string  `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	AvailableSeats int32   `json:"available_seats" validate:"required,min=1"`
	PricePerSeat   float64 `json:"price_per_seat" validate:"required,min=0"`
	VenueID        int32   `json:"venue_id" validate:"required"`
}

// UpdateShowtimeRequest represents the request to update a showtime.
type UpdateShowtimeRequest struct {
	StartTime      *string  `json:"start_time" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime        *string  `json:"end_time" validate:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	AvailableSeats *int32   `json:"available_seats" validate:"omitempty,min=0"`
	PricePerSeat   *float64 `json:"price_per_seat" validate:"omitempty,min=0"`
	VenueID        *int32   `json:"venue_id"`
}
