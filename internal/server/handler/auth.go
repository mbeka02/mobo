package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/internal/model"
	customMiddleware "github.com/mbeka02/ticketing-service/internal/server/middleware"
	"github.com/mbeka02/ticketing-service/internal/server/service"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService      service.AuthService
	tokenMaker       auth.Maker
	isProduction     bool
	accessDuration   time.Duration
	refreshDuration  time.Duration
}

func NewAuthHandler(svc service.AuthService, maker auth.Maker, isProduction bool, accessDuration, refreshDuration time.Duration) *AuthHandler {
	return &AuthHandler{
		authService:      svc,
		tokenMaker:       maker,
		isProduction:     isProduction,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateLocalUserRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		logger.WarnCtx(ctx, "invalid signup request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.authService.RegisterLocalUser(ctx, req.Email, req.Fullname, req.Password, req.TelephoneNumber)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			respondWithError(w, http.StatusConflict, err)
			return
		}
		logger.ErrorCtx(ctx, "failed to register user", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Set JWT token cookies
	if err := auth.SetTokenCookies(w, h.tokenMaker, user.ID, user.Email, h.isProduction, h.accessDuration, h.refreshDuration); err != nil {
		logger.ErrorCtx(ctx, "failed to set token cookies", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Status:  http.StatusCreated,
		Message: "user registered successfully",
		Data:    user.ToResponse(),
	})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.LoginRequest
	if err := parseAndValidateRequest(r, &req); err != nil {
		logger.WarnCtx(ctx, "invalid login request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.authService.LoginLocalUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondWithError(w, http.StatusUnauthorized, err)
			return
		}
		if errors.Is(err, service.ErrOAuthOnlyAccount) {
			respondWithError(w, http.StatusConflict, err)
			return
		}
		logger.ErrorCtx(ctx, "login failed", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Set JWT token cookies
	if err := auth.SetTokenCookies(w, h.tokenMaker, user.ID, user.Email, h.isProduction, h.accessDuration, h.refreshDuration); err != nil {
		logger.ErrorCtx(ctx, "failed to set token cookies", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "login successful",
		Data:    user.ToResponse(),
	})
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	auth.ClearTokenCookies(w)

	logger.InfoCtx(ctx, "user logged out")

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "logged out successfully",
	})
}

func (h *AuthHandler) BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := logger.WithRequestID(r.Context(), requestID)
	r = r.WithContext(ctx)

	provider := chi.URLParam(r, "provider")
	if provider == "" {
		logger.WarnCtx(ctx, "auth attempt with missing provider")
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}

	r = r.WithContext(context.WithValue(ctx, "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) GetAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	ctx := logger.WithRequestID(r.Context(), requestID)
	r = r.WithContext(ctx)

	provider := chi.URLParam(r, "provider")
	if provider == "" {
		logger.WarnCtx(ctx, "auth callback with invalid provider")
		http.Redirect(w, r, "http://localhost:5173/login?error=invalid_provider", http.StatusFound)
		return
	}

	logger.InfoCtx(ctx, "processing OAuth callback",
		zap.String("provider", provider),
	)

	r = r.WithContext(context.WithValue(ctx, "provider", provider))
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		logger.ErrorCtx(ctx, "OAuth authentication failed",
			zap.Error(err),
			zap.String("provider", provider),
		)
		http.Redirect(w, r, "http://localhost:5173/login?error=auth_failed", http.StatusFound)
		return
	}
	// Convert to service DTO
	oauthUser := service.OAuthUserData{
		Email:          gothUser.Email,
		FirstName:      gothUser.FirstName,
		LastName:       gothUser.LastName,
		Provider:       gothUser.Provider,
		ProviderUserID: gothUser.UserID,
		AvatarURL:      gothUser.AvatarURL,
	}

	user, err := h.authService.CreateOrLoginOAuthUser(ctx, oauthUser)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create or login OAuth user",
			zap.Error(err),
			zap.String("provider", provider),
			zap.String("email", gothUser.Email),
		)
		http.Redirect(w, r, "http://localhost:5173/login?error=account_creation_failed", http.StatusFound)
		return
	}

	// Set JWT token cookies for unified auth
	if err := auth.SetTokenCookies(w, h.tokenMaker, user.ID, user.Email, h.isProduction, h.accessDuration, h.refreshDuration); err != nil {
		logger.ErrorCtx(ctx, "failed to set token cookies after OAuth",
			zap.Error(err),
			zap.String("user_id", user.ID.String()),
		)
		http.Redirect(w, r, "http://localhost:5173/login?error=token_failed", http.StatusFound)
		return
	}

	http.Redirect(w, r, "http://localhost:5173/home", http.StatusFound)
}

func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := customMiddleware.UserIDFromContext(ctx)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	email, _ := ctx.Value(customMiddleware.EmailKey).(string)

	logger.DebugCtx(ctx, "fetching current user",
		zap.String("user_id", userID.String()),
	)

	respondWithJSON(w, http.StatusOK, APIResponse{
		Status:  http.StatusOK,
		Message: "current user",
		Data: map[string]string{
			"user_id": userID.String(),
			"email":   email,
		},
	})
}
