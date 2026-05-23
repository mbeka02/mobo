package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
)

var (
	ErrInvalidShowtimeRange = errors.New("start time must be before end time")
)

type ShowtimeService interface {
	CreateShowtime(ctx context.Context, req model.CreateShowtimeRequest) (*model.Showtime, error)
	GetShowtime(ctx context.Context, id int64) (*model.Showtime, error)
	ListShowtimesByMovie(ctx context.Context, movieId int64) ([]model.Showtime, error)
	ListShowtimesAdmin(ctx context.Context, limit, offset int32) ([]model.Showtime, error)
	UpdateShowtime(ctx context.Context, id int64, req model.UpdateShowtimeRequest) (*model.Showtime, error)
	DeleteShowtime(ctx context.Context, id int64) error
}

type showtimeService struct {
	repo repository.ShowtimeRepository
}

func NewShowtimeService(repo repository.ShowtimeRepository) ShowtimeService {
	return &showtimeService{repo}
}

func (s *showtimeService) CreateShowtime(ctx context.Context, req model.CreateShowtimeRequest) (*model.Showtime, error) {
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, err
	}

	if !end.After(start) {
		return nil, ErrInvalidShowtimeRange
	}

	// We use strings for numeric to ensure precision
	price := pgtype.Numeric{}
	price.Scan(req.PricePerSeat)

	params := database.CreateShowtimeParams{
		MovieID:        req.MovieID,
		StartTime:      start,
		EndTime:        end,
		AvailableSeats: req.AvailableSeats,
		PricePerSeat:   price,
		VenueID:        req.VenueID,
	}
	return s.repo.CreateShowtime(ctx, params)
}

func (s *showtimeService) GetShowtime(ctx context.Context, id int64) (*model.Showtime, error) {
	return s.repo.GetShowtimeById(ctx, id)
}

func (s *showtimeService) ListShowtimesByMovie(ctx context.Context, movieId int64) ([]model.Showtime, error) {
	return s.repo.GetShowtimesByMovie(ctx, movieId)
}

func (s *showtimeService) ListShowtimesAdmin(ctx context.Context, limit, offset int32) ([]model.Showtime, error) {
	return s.repo.GetShowtimesAdmin(ctx, limit, offset)
}

func (s *showtimeService) UpdateShowtime(ctx context.Context, id int64, req model.UpdateShowtimeRequest) (*model.Showtime, error) {
	params := database.UpdateShowtimeParams{
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

	return s.repo.UpdateShowtime(ctx, params)
}

func (s *showtimeService) DeleteShowtime(ctx context.Context, id int64) error {
	return s.repo.DeleteShowtime(ctx, id)
}
