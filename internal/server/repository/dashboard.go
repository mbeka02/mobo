package repository

import (
	"context"

	"github.com/mbeka02/ticketing-service/internal/database"
)

type DashboardRepository interface {
	GetDashboardStats(ctx context.Context) (database.GetDashboardStatsRow, error)
	GetMonthlyRevenue(ctx context.Context) ([]database.GetMonthlyRevenueRow, error)
}

type dashboardRepository struct {
	store *database.Store
}

func NewDashboardRepository(store *database.Store) DashboardRepository {
	return &dashboardRepository{store}
}

func (r *dashboardRepository) GetDashboardStats(ctx context.Context) (database.GetDashboardStatsRow, error) {
	return r.store.GetDashboardStats(ctx)
}

func (r *dashboardRepository) GetMonthlyRevenue(ctx context.Context) ([]database.GetMonthlyRevenueRow, error) {
	return r.store.GetMonthlyRevenue(ctx)
}
