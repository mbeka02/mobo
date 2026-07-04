package movie

import "time"

// Movie represents a movie in the system.
type Movie struct {
	ID          int64
	Title       string
	Description string
	Runtime     int32
	Genre       string
	AgeRating   string
	Director    string
	PosterUrl   string
	ReleaseDate time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// ToResponse converts a Movie to a MovieResponse.
func (m *Movie) ToResponse() MovieResponse {
	var updatedAt time.Time
	if m.UpdatedAt != nil {
		updatedAt = *m.UpdatedAt
	}

	return MovieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Runtime:     m.Runtime,
		Genre:       m.Genre,
		AgeRating:   m.AgeRating,
		Director:    m.Director,
		PosterUrl:   m.PosterUrl,
		ReleaseDate: m.ReleaseDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   updatedAt,
	}
}

// MovieResponse represents the API response for a movie.
type MovieResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Runtime     int32     `json:"runtime"`
	Genre       string    `json:"genre"`
	AgeRating   string    `json:"age_rating"`
	Director    string    `json:"director"`
	PosterUrl   string    `json:"poster_url"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// AddMovieRequest represents the request to add a new movie.
type AddMovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Runtime     int32  `json:"runtime" validate:"required,min=1"`
	Genre       string `json:"genre" validate:"required"`
	AgeRating   string `json:"age_rating" validate:"required"`
	Director    string `json:"director" validate:"required"`
	PosterUrl   string `json:"poster_url" validate:"required,url"`
	ReleaseDate string `json:"release_date" validate:"required,datetime=2006-01-02"`
}

// UpdateMovieRequest represents the request to update a movie.
type UpdateMovieRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Runtime     *int32  `json:"runtime" validate:"omitempty,min=1"`
	Genre       *string `json:"genre"`
	AgeRating   *string `json:"age_rating"`
	Director    *string `json:"director"`
	PosterUrl   *string `json:"poster_url" validate:"omitempty,url"`
	ReleaseDate *string `json:"release_date" validate:"omitempty,datetime=2006-01-02"`
}
