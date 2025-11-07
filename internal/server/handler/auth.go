package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/mbeka02/ticketing-service/internal/server/service"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
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

	session, _ := gothic.Store.Get(r, "user-session")
	session.Values["user_id"] = user.ID.String()
	session.Values["email"] = user.Email
	err = session.Save(r, w)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to save user session",
			zap.Error(err),
			zap.String("user_id", user.ID.String()),
		)
		http.Redirect(w, r, "http://localhost:5173/login?error=session_failed", http.StatusFound)
		return
	}
	http.Redirect(w, r, "http://localhost:5173/home", http.StatusFound)
}
