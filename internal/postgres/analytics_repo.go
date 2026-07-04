package postgres

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/analytics"
)

type analyticsRepo struct {
	store *Store
}

// NewAnalyticsRepository creates a new postgres analytics repository.
func NewAnalyticsRepository(store *Store) analytics.Repository {
	return &analyticsRepo{store}
}

func (r *analyticsRepo) GetStats(ctx context.Context) (analytics.StatsRow, error) {
	row, err := r.store.GetDashboardStats(ctx)
	if err != nil {
		return analytics.StatsRow{}, err
	}
	return analytics.StatsRow{
		TotalRevenue: row.TotalRevenue,
		TicketsSold:  row.TicketsSold,
		ActiveVenues: row.ActiveVenues,
		ActiveMovies: row.ActiveMovies,
	}, nil
}

func (r *analyticsRepo) GetMonthlyRevenue(ctx context.Context) ([]analytics.MonthlyRevenueRow, error) {
	rows, err := r.store.GetMonthlyRevenue(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]analytics.MonthlyRevenueRow, 0, len(rows))
	for _, row := range rows {
		entry := analytics.MonthlyRevenueRow{
			Revenue: row.Revenue,
		}
		if row.Month.Valid {
			entry.Month = row.Month.Time
			entry.Valid = true
		}
		result = append(result, entry)
	}
	return result, nil
}
