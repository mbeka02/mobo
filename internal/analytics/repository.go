package analytics

import "context"

// StatsRow represents the raw data from the repository for dashboard stats.
type StatsRow struct {
	TotalRevenue string
	TicketsSold  int64
	ActiveVenues int32
	ActiveMovies int32
}

// Repository defines the data access contract for the analytics domain.
type Repository interface {
	GetStats(ctx context.Context) (StatsRow, error)
	GetMonthlyRevenue(ctx context.Context) ([]MonthlyRevenueRow, error)
}
