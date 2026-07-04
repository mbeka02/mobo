package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mbeka02/ticketing-service/internal/venue"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

// VenueHandler handles HTTP requests for the venue domain.
type VenueHandler struct {
	svc venue.Service
}

// NewVenueHandler creates a new VenueHandler.
func NewVenueHandler(svc venue.Service) *VenueHandler {
	return &VenueHandler{svc: svc}
}

func (h *VenueHandler) CreateVenueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req venue.CreateVenueRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	v, err := h.svc.CreateVenue(ctx, req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create venue", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Status:  http.StatusCreated,
		Message: "venue created successfully",
		Data:    v.ToResponse(),
	})
}

func (h *VenueHandler) GetVenueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "venueId")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	v, err := h.svc.GetVenue(ctx, int32(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    v.ToResponse(),
	})
}

func (h *VenueHandler) ListVenuesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	venues, err := h.svc.ListVenues(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to list venues", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var res []venue.VenueResponse
	for _, v := range venues {
		res = append(res, v.ToResponse())
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}
