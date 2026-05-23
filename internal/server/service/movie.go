package service

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mbeka02/ticketing-service/internal/database"
	"github.com/mbeka02/ticketing-service/internal/model"
	"github.com/mbeka02/ticketing-service/internal/server/repository"
)

type MovieService interface {
	AddMovie(ctx context.Context, req model.AddMovieRequest) (*model.Movie, error)
	GetMovie(ctx context.Context, id int64) (*model.Movie, error)
	ListMoviesAdmin(ctx context.Context, limit, offset int32) ([]model.Movie, error)
	ListMoviesPublic(ctx context.Context, limit, offset int32) ([]model.Movie, error)
	UpdateMovie(ctx context.Context, id int64, req model.UpdateMovieRequest) (*model.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
}

type movieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) MovieService {
	return &movieService{repo}
}

func (s *movieService) AddMovie(ctx context.Context, req model.AddMovieRequest) (*model.Movie, error) {
	params, err := req.ToParams()
	if err != nil {
		return nil, err
	}
	return s.repo.AddMovie(ctx, params)
}

func (s *movieService) GetMovie(ctx context.Context, id int64) (*model.Movie, error) {
	return s.repo.GetMovieById(ctx, id)
}

func (s *movieService) ListMoviesAdmin(ctx context.Context, limit, offset int32) ([]model.Movie, error) {
	return s.repo.GetMoviesAdmin(ctx, limit, offset)
}

func (s *movieService) ListMoviesPublic(ctx context.Context, limit, offset int32) ([]model.Movie, error) {
	return s.repo.GetMoviesPublic(ctx, limit, offset)
}

func (s *movieService) UpdateMovie(ctx context.Context, id int64, req model.UpdateMovieRequest) (*model.Movie, error) {
	params := database.UpdateMovieParams{
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

	return s.repo.UpdateMovie(ctx, params)
}

func (s *movieService) DeleteMovie(ctx context.Context, id int64) error {
	return s.repo.DeleteMovie(ctx, id)
}
