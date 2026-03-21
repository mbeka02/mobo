package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/mbeka02/ticketing-service/internal/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

func AuthMiddleware(maker auth.Maker, isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID uuid.UUID
			var email string
			var found bool

			// 1. Try access token
			if cookie, err := r.Cookie(auth.AccessTokenCookie); err == nil {
				if claims, err := maker.Verify(cookie.Value); err == nil {
					userID = claims.UserID
					email = claims.Email
					found = true
				} else if errors.Is(err, auth.ErrExpiredToken) {
					// 2. Access token expired — try refresh token
					if refresh, err := r.Cookie(auth.RefreshTokenCookie); err == nil {
						if claims, err := maker.Verify(refresh.Value); err == nil {
							// Silently rotate both cookies
							if err := auth.SetTokenCookies(w, maker, claims.UserID, claims.Email, isProd); err != nil {
								http.Error(w, "unauthorized", http.StatusUnauthorized)
								return
							}
							userID = claims.UserID
							email = claims.Email
							found = true
						}
					}
				}
			}

			// 3. Try Gothic session (OAuth users)
			if !found {
				if session, err := gothic.Store.Get(r, "user-session"); err == nil {
					if id, ok := session.Values["user_id"].(string); ok && id != "" {
						if parsedID, err := uuid.Parse(id); err == nil {
							userID = parsedID
							found = true
						}
					}
				}
			}

			if !found {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, EmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper for handlers to pull the user ID back out
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
