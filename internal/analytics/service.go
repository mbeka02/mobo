package analytics

import (
	"context"
	"time"
)

// Service defines the business operations for the analytics domain.
type Service interface {
	GetStats(ctx context.Context) (*StatsResponse, error)
	GetMonthlyRevenue(ctx context.Context) ([]MonthlyRevenueEntry, error)
}

type service struct {
	repo Repository
}

// NewService creates a new analytics service.
func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetStats(ctx context.Context) (*StatsResponse, error) {
	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	return &StatsResponse{
		TotalRevenue: stats.TotalRevenue,
		TicketsSold:  stats.TicketsSold,
		ActiveVenues: stats.ActiveVenues,
		ActiveMovies: stats.ActiveMovies,
	}, nil
}

func (s *service) GetMonthlyRevenue(ctx context.Context) ([]MonthlyRevenueEntry, error) {
	rows, err := s.repo.GetMonthlyRevenue(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]MonthlyRevenueEntry, 0, len(rows))
	for _, r := range rows {
		var monthStr string
		if r.Valid {
			monthStr = r.Month.Format(time.DateOnly)
		}
		entries = append(entries, MonthlyRevenueEntry{
			Month:   monthStr,
			Revenue: r.Revenue,
		})
	}
	return entries, nil
}
