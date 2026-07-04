package api

import (
	"net/http"

	"github.com/mbeka02/ticketing-service/internal/analytics"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

// AnalyticsHandler handles HTTP requests for the analytics domain.
type AnalyticsHandler struct {
	svc analytics.Service
}

// NewAnalyticsHandler creates a new AnalyticsHandler.
func NewAnalyticsHandler(svc analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) GetDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.svc.GetStats(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to get dashboard stats", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	revenue, err := h.svc.GetMonthlyRevenue(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to get monthly revenue", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"stats":           stats,
			"monthly_revenue": revenue,
		},
	})
}
