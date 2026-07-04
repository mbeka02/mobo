package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/dbgen"
	"github.com/mbeka02/ticketing-service/internal/movie"
)

type movieRepo struct {
	store *Store
}

// NewMovieRepository creates a new postgres movie repository.
func NewMovieRepository(store *Store) movie.Repository {
	return &movieRepo{store}
}

func (r *movieRepo) Add(ctx context.Context, req movie.AddMovieRequest) (*movie.Movie, error) {
	parsedDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return nil, err
	}

	dbMovie, err := r.store.AddMovie(ctx, dbgen.AddMovieParams{
		Title:       req.Title,
		Description: req.Description,
		Runtime:     req.Runtime,
		Genre:       req.Genre,
		AgeRating:   req.AgeRating,
		Director:    req.Director,
		PosterUrl:   req.PosterUrl,
		ReleaseDate: pgtype.Date{Time: parsedDate, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return fromDatabaseMovie(&dbMovie), nil
}

func (r *movieRepo) GetByID(ctx context.Context, id int64) (*movie.Movie, error) {
	row, err := r.store.GetMovieById(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromDatabaseGetMovieByIdRow(&row), nil
}

func (r *movieRepo) ListAdmin(ctx context.Context, limit, offset int32) ([]movie.Movie, error) {
	movies, err := r.store.GetMoviesAdmin(ctx, dbgen.GetMoviesAdminParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]movie.Movie, 0, len(movies))
	for _, m := range movies {
		res = append(res, *fromDatabaseGetMoviesAdminRow(&m))
	}
	return res, nil
}

func (r *movieRepo) ListPublic(ctx context.Context, limit, offset int32) ([]movie.Movie, error) {
	movies, err := r.store.GetMoviesPublic(ctx, dbgen.GetMoviesPublicParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]movie.Movie, 0, len(movies))
	for _, m := range movies {
		res = append(res, *fromDatabaseMovie(&m))
	}
	return res, nil
}

func (r *movieRepo) Update(ctx context.Context, id int64, req movie.UpdateMovieRequest) (*movie.Movie, error) {
	params := dbgen.UpdateMovieParams{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Runtime:     req.Runtime,
		Genre:       req.Genre,
		AgeRating:   req.AgeRating,
		Director:    req.Director,
		PosterUrl:   req.PosterUrl,
	}

	if req.ReleaseDate != nil {
		if t, err := time.Parse("2006-01-02", *req.ReleaseDate); err == nil {
			params.ReleaseDate = pgtype.Date{Time: t, Valid: true}
		}
	}

	dbMovie, err := r.store.UpdateMovie(ctx, params)
	if err != nil {
		return nil, err
	}
	return fromDatabaseMovie(&dbMovie), nil
}

func (r *movieRepo) Delete(ctx context.Context, id int64) error {
	return r.store.DeleteMovie(ctx, id)
}

// Conversion helpers

func fromDatabaseMovie(dbMovie *dbgen.Movie) *movie.Movie {
	var updatedAt *time.Time
	if dbMovie.UpdatedAt.Valid {
		updatedAt = &dbMovie.UpdatedAt.Time
	}
	var releaseDate time.Time
	if dbMovie.ReleaseDate.Valid {
		releaseDate = dbMovie.ReleaseDate.Time
	}

	return &movie.Movie{
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

func fromDatabaseGetMoviesAdminRow(row *dbgen.GetMoviesAdminRow) *movie.Movie {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}
	var releaseDate time.Time
	if row.ReleaseDate.Valid {
		releaseDate = row.ReleaseDate.Time
	}

	return &movie.Movie{
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

func fromDatabaseGetMovieByIdRow(row *dbgen.GetMovieByIdRow) *movie.Movie {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}
	var releaseDate time.Time
	if row.ReleaseDate.Valid {
		releaseDate = row.ReleaseDate.Time
	}

	return &movie.Movie{
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
