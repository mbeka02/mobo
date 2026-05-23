package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/database"
)

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

func FromDatabaseMovie(dbMovie *database.Movie) *Movie {
	var updatedAt *time.Time
	if dbMovie.UpdatedAt.Valid {
		updatedAt = &dbMovie.UpdatedAt.Time
	}

	var releaseDate time.Time
	if dbMovie.ReleaseDate.Valid {
		releaseDate = dbMovie.ReleaseDate.Time
	}

	return &Movie{
		ID:          dbMovie.ID,
		Title:       dbMovie.Title,
		Description: dbMovie.Description,
		Runtime:     dbMovie.Runtime,
		Genre:       dbMovie.Genre,
		AgeRating:   dbMovie.AgeRating,
		Director:    dbMovie.Director,
		PosterUrl:   dbMovie.PosterUrl,
		ReleaseDate: releaseDate,
		CreatedAt:   dbMovie.CreatedAt,
		UpdatedAt:   updatedAt,
	}
}

func FromDatabaseGetMoviesAdminRow(row *database.GetMoviesAdminRow) *Movie {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}

	var releaseDate time.Time
	if row.ReleaseDate.Valid {
		releaseDate = row.ReleaseDate.Time
	}

	return &Movie{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Runtime:     row.Runtime,
		Genre:       row.Genre,
		AgeRating:   row.AgeRating,
		Director:    row.Director,
		PosterUrl:   row.PosterUrl,
		ReleaseDate: releaseDate,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   updatedAt,
	}
}

func FromDatabaseGetMovieByIdRow(row *database.GetMovieByIdRow) *Movie {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}

	var releaseDate time.Time
	if row.ReleaseDate.Valid {
		releaseDate = row.ReleaseDate.Time
	}

	return &Movie{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Runtime:     row.Runtime,
		Genre:       row.Genre,
		AgeRating:   row.AgeRating,
		Director:    row.Director,
		PosterUrl:   row.PosterUrl,
		ReleaseDate: releaseDate,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   updatedAt,
	}
}

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

func (req *AddMovieRequest) ToParams() (database.AddMovieParams, error) {
	parsedDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return database.AddMovieParams{}, err
	}
	return database.AddMovieParams{
		Title:       req.Title,
		Description: req.Description,
		Runtime:     req.Runtime,
		Genre:       req.Genre,
		AgeRating:   req.AgeRating,
		Director:    req.Director,
		PosterUrl:   req.PosterUrl,
		ReleaseDate: pgtype.Date{Time: parsedDate, Valid: true},
	}, nil
}

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
