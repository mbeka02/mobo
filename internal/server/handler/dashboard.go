package handler

import (
	"net/http"

	"github.com/mbeka02/ticketing-service/internal/server/service"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
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
