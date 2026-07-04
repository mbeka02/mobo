package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/dbgen"
	"github.com/mbeka02/ticketing-service/internal/showtime"
)

type showtimeRepo struct {
	store *Store
}

// NewShowtimeRepository creates a new postgres showtime repository.
func NewShowtimeRepository(store *Store) showtime.Repository {
	return &showtimeRepo{store}
}

func (r *showtimeRepo) Create(ctx context.Context, req showtime.CreateShowtimeRequest) (*showtime.Showtime, error) {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	price := pgtype.Numeric{}
	price.Scan(req.PricePerSeat)

	dbShowtime, err := r.store.CreateShowtime(ctx, dbgen.CreateShowtimeParams{
		MovieID:        req.MovieID,
		StartTime:      start,
		EndTime:        end,
		AvailableSeats: req.AvailableSeats,
		PricePerSeat:   price,
		VenueID:        req.VenueID,
	})
	if err != nil {
		return nil, err
	}
	return fromDatabaseShowtime(&dbShowtime), nil
}

func (r *showtimeRepo) GetByID(ctx context.Context, id int64) (*showtime.Showtime, error) {
	dbShowtime, err := r.store.GetShowtimeById(ctx, id)
	if err != nil {
		return nil, err
	}
	return fromDatabaseShowtime(&dbShowtime), nil
}

func (r *showtimeRepo) ListByMovie(ctx context.Context, movieId int64) ([]showtime.Showtime, error) {
	rows, err := r.store.GetShowtimesByMovie(ctx, movieId)
	if err != nil {
		return nil, err
	}

	res := make([]showtime.Showtime, 0, len(rows))
	for _, s := range rows {
		res = append(res, *fromDatabaseGetShowtimesByMovieRow(&s))
	}
	return res, nil
}

func (r *showtimeRepo) ListAdmin(ctx context.Context, limit, offset int32) ([]showtime.Showtime, error) {
	rows, err := r.store.GetShowtimesAdmin(ctx, dbgen.GetShowtimesAdminParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]showtime.Showtime, 0, len(rows))
	for _, s := range rows {
		res = append(res, *fromDatabaseGetShowtimesAdminRow(&s))
	}
	return res, nil
}

func (r *showtimeRepo) Update(ctx context.Context, id int64, req showtime.UpdateShowtimeRequest) (*showtime.Showtime, error) {
	params := dbgen.UpdateShowtimeParams{
		ID:             id,
		AvailableSeats: req.AvailableSeats,
		VenueID:        req.VenueID,
	}

	if req.StartTime != nil {
		if t, err := time.Parse(time.RFC3339, *req.StartTime); err == nil {
			params.StartTime = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	if req.EndTime != nil {
		if t, err := time.Parse(time.RFC3339, *req.EndTime); err == nil {
			params.EndTime = pgtype.Timestamptz{Time: t, Valid: true}
		}
	}
	if req.PricePerSeat != nil {
		price := pgtype.Numeric{}
		price.Scan(*req.PricePerSeat)
		params.PricePerSeat = price
	}

	dbShowtime, err := r.store.UpdateShowtime(ctx, params)
	if err != nil {
		return nil, err
	}
	return fromDatabaseShowtime(&dbShowtime), nil
}

func (r *showtimeRepo) Delete(ctx context.Context, id int64) error {
	return r.store.DeleteShowtime(ctx, id)
}

// Conversion helpers

func fromDatabaseShowtime(dbShowtime *dbgen.Showtime) *showtime.Showtime {
	var updatedAt *time.Time
	if dbShowtime.UpdatedAt.Valid {
		updatedAt = &dbShowtime.UpdatedAt.Time
	}
	price, _ := dbShowtime.PricePerSeat.Float64Value()

	return &showtime.Showtime{
		ID:             dbShowtime.ID,
		MovieID:        dbShowtime.MovieID,
		StartTime:      dbShowtime.StartTime,
		EndTime:        dbShowtime.EndTime,
		AvailableSeats: dbShowtime.AvailableSeats,
		PricePerSeat:   price.Float64,
		VenueID:        dbShowtime.VenueID,
		CreatedAt:      dbShowtime.CreatedAt,
		UpdatedAt:      updatedAt,
	}
}

func fromDatabaseGetShowtimesByMovieRow(row *dbgen.GetShowtimesByMovieRow) *showtime.Showtime {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}
	price, _ := row.PricePerSeat.Float64Value()

	return &showtime.Showtime{
		ID:             row.ID,
		MovieID:        row.MovieID,
		StartTime:      row.StartTime,
		EndTime:        row.EndTime,
		AvailableSeats: row.AvailableSeats,
		PricePerSeat:   price.Float64,
		VenueID:        row.VenueID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      updatedAt,
		VenueName:      &row.VenueName,
		VenueCity:      &row.VenueCity,
	}
}

func fromDatabaseGetShowtimesAdminRow(row *dbgen.GetShowtimesAdminRow) *showtime.Showtime {
	var updatedAt *time.Time
	if row.UpdatedAt.Valid {
		updatedAt = &row.UpdatedAt.Time
	}
	price, _ := row.PricePerSeat.Float64Value()

	return &showtime.Showtime{
		ID:             row.ID,
		MovieID:        row.MovieID,
		StartTime:      row.StartTime,
		EndTime:        row.EndTime,
		AvailableSeats: row.AvailableSeats,
		PricePerSeat:   price.Float64,
		VenueID:        row.VenueID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      updatedAt,
		MovieTitle:     &row.MovieTitle,
		VenueName:      &row.VenueName,
	}
}
