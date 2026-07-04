package analytics

import "time"

// StatsResponse contains the aggregated dashboard statistics.
type StatsResponse struct {
	TotalRevenue string `json:"total_revenue"`
	TicketsSold  int64  `json:"tickets_sold"`
	ActiveVenues int32  `json:"active_venues"`
	ActiveMovies int32  `json:"active_movies"`
}

// MonthlyRevenueEntry represents revenue data for a single month.
type MonthlyRevenueEntry struct {
	Month   string `json:"month"`
	Revenue string `json:"revenue"`
}

// MonthlyRevenueRow represents the raw data from the repository for monthly revenue.
type MonthlyRevenueRow struct {
	Month   time.Time
	Valid   bool
	Revenue string
}
