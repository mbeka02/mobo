package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
	"github.com/mbeka02/ticketing-service/internal/server/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) GetAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {

		http.Redirect(w, r, "http://localhost:5173/login?error=invalid_provider", http.StatusFound)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
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

	user, err := h.authService.CreateOrLoginOAuthUser(r.Context(), oauthUser)
	if err != nil {
		http.Redirect(w, r, "http://localhost:5173/login?error=account_creation_failed", http.StatusFound)
		return
	}

	session, _ := gothic.Store.Get(r, "user-session")
	session.Values["user_id"] = user.ID.String()
	session.Values["email"] = user.Email
	err = session.Save(r, w)
	if err != nil {
		http.Redirect(w, r, "http://localhost:5173/login?error=session_failed", http.StatusFound)
		return
	}

	http.Redirect(w, r, "http://localhost:5173/home", http.StatusFound)
}
