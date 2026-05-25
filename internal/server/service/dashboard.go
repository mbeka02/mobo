package service

import (
	"context"
	"time"

	"github.com/mbeka02/ticketing-service/internal/server/repository"
)

type DashboardStatsResponse struct {
	TotalRevenue string `json:"total_revenue"`
	TicketsSold  int64  `json:"tickets_sold"`
	ActiveVenues int32  `json:"active_venues"`
	ActiveMovies int32  `json:"active_movies"`
}

type MonthlyRevenueEntry struct {
	Month   string `json:"month"`
	Revenue string `json:"revenue"`
}

type DashboardService interface {
	GetStats(ctx context.Context) (*DashboardStatsResponse, error)
	GetMonthlyRevenue(ctx context.Context) ([]MonthlyRevenueEntry, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo}
}

func (s *dashboardService) GetStats(ctx context.Context) (*DashboardStatsResponse, error) {
	stats, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	return &DashboardStatsResponse{
		TotalRevenue: stats.TotalRevenue,
		TicketsSold:  stats.TicketsSold,
		ActiveVenues: stats.ActiveVenues,
		ActiveMovies: stats.ActiveMovies,
	}, nil
}

func (s *dashboardService) GetMonthlyRevenue(ctx context.Context) ([]MonthlyRevenueEntry, error) {
	rows, err := s.repo.GetMonthlyRevenue(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]MonthlyRevenueEntry, 0, len(rows))
	for _, r := range rows {
		var monthStr string
		if r.Month.Valid {
			monthStr = r.Month.Time.Format(time.DateOnly)
		}
		entries = append(entries, MonthlyRevenueEntry{
			Month:   monthStr,
			Revenue: r.Revenue,
		})
	}
	return entries, nil
}
