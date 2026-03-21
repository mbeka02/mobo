package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/mbeka02/ticketing-service/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(maker auth.Maker, isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID uuid.UUID
			var found bool
			// FIXME: No token rotation exists at the moment
			// Try JWT access cookie
			if cookie, err := r.Cookie(auth.AccessTokenCookie); err == nil {
				if claims, err := maker.Verify(cookie.Value); err == nil {
					userID = claims.UserID
					found = true
				}
			}

			// Try Gothic session (OIDC)
			if !found {
				if session, err := gothic.Store.Get(r, "user-session"); err == nil {
					if id, ok := session.Values["user_id"].(string); ok && id != "" {
						if parsed, err := uuid.Parse(id); err == nil {
							userID = parsed
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
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper for handlers to pull the user ID back out
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}
