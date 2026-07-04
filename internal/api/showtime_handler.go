package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mbeka02/ticketing-service/internal/showtime"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

// ShowtimeHandler handles HTTP requests for the showtime domain.
type ShowtimeHandler struct {
	svc showtime.Service
}

// NewShowtimeHandler creates a new ShowtimeHandler.
func NewShowtimeHandler(svc showtime.Service) *ShowtimeHandler {
	return &ShowtimeHandler{svc: svc}
}

func (h *ShowtimeHandler) CreateShowtimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req showtime.CreateShowtimeRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	s, err := h.svc.CreateShowtime(ctx, req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create showtime", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Status:  http.StatusCreated,
		Message: "showtime created successfully",
		Data:    s.ToResponse(),
	})
}

func (h *ShowtimeHandler) GetShowtimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "showtimeId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	s, err := h.svc.GetShowtime(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    s.ToResponse(),
	})
}

func (h *ShowtimeHandler) ListShowtimesByMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "movieId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	showtimes, err := h.svc.ListShowtimesByMovie(ctx, id)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to list showtimes by movie", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var res []showtime.ShowtimeResponse
	for _, s := range showtimes {
		res = append(res, s.ToResponse())
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h *ShowtimeHandler) ListShowtimesAdminHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset := parsePagination(r)

	showtimes, err := h.svc.ListShowtimesAdmin(ctx, limit, offset)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to list admin showtimes", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var res []showtime.ShowtimeResponse
	for _, s := range showtimes {
		res = append(res, s.ToResponse())
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h *ShowtimeHandler) UpdateShowtimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "showtimeId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var req showtime.UpdateShowtimeRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	s, err := h.svc.UpdateShowtime(ctx, id, req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to update showtime", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "showtime updated successfully",
		Data:    s.ToResponse(),
	})
}

func (h *ShowtimeHandler) DeleteShowtimeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "showtimeId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeleteShowtime(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "failed to delete showtime", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "showtime deleted successfully",
	})
}
