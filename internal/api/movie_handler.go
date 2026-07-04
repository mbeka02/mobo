package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mbeka02/ticketing-service/internal/movie"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

// MovieHandler handles HTTP requests for the movie domain.
type MovieHandler struct {
	svc movie.Service
}

// NewMovieHandler creates a new MovieHandler.
func NewMovieHandler(svc movie.Service) *MovieHandler {
	return &MovieHandler{svc: svc}
}

func (h *MovieHandler) AddMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req movie.AddMovieRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	m, err := h.svc.AddMovie(ctx, req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to add movie", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Status:  http.StatusCreated,
		Message: "movie created successfully",
		Data:    m.ToResponse(),
	})
}

func (h *MovieHandler) GetMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "movieId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	m, err := h.svc.GetMovie(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    m.ToResponse(),
	})
}

func (h *MovieHandler) ListMoviesAdminHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset := parsePagination(r)

	movies, err := h.svc.ListMoviesAdmin(ctx, limit, offset)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to list admin movies", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var res []movie.MovieResponse
	for _, m := range movies {
		res = append(res, m.ToResponse())
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h *MovieHandler) ListMoviesPublicHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset := parsePagination(r)

	movies, err := h.svc.ListMoviesPublic(ctx, limit, offset)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to list public movies", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	var res []movie.MovieResponse
	for _, m := range movies {
		res = append(res, m.ToResponse())
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    res,
	})
}

func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "movieId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var req movie.UpdateMovieRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	m, err := h.svc.UpdateMovie(ctx, id, req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to update movie", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "movie updated successfully",
		Data:    m.ToResponse(),
	})
}

func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "movieId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.DeleteMovie(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "failed to delete movie", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "movie deleted successfully",
	})
}
