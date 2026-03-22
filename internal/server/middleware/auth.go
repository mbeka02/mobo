package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/ticketing-service/internal/auth"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

func AuthMiddleware(maker auth.Maker, isProd bool, accessDuration, refreshDuration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var userID uuid.UUID
			var email string
			var found bool

			// 1. Try access token cookie
			if cookie, err := r.Cookie(auth.AccessTokenCookie); err == nil {
				claims, err := maker.Verify(cookie.Value)
				if err == nil {
					// Verify this is an access token, not a refresh token
					if claims.TokenType != auth.AccessToken {
						logger.WarnCtx(ctx, "non-access token used in access cookie")
						http.Error(w, "unauthorized", http.StatusUnauthorized)
						return
					}
					userID = claims.UserID
					email = claims.Email
					found = true
					logger.DebugCtx(ctx, "authenticated via access token",
						zap.String("user_id", userID.String()),
					)
				} else if errors.Is(err, auth.ErrExpiredToken) {
					// 2. Access token expired — try refresh token rotation
					logger.DebugCtx(ctx, "access token expired, attempting refresh")
					if refresh, err := r.Cookie(auth.RefreshTokenCookie); err == nil {
						refreshClaims, err := maker.Verify(refresh.Value)
						if err == nil {
							// Verify this is a refresh token
							if refreshClaims.TokenType != auth.RefreshToken {
								logger.WarnCtx(ctx, "non-refresh token used in refresh cookie")
								http.Error(w, "unauthorized", http.StatusUnauthorized)
								return
							}
							// Silently rotate both cookies
							if err := auth.SetTokenCookies(w, maker, refreshClaims.UserID, refreshClaims.Email, isProd, accessDuration, refreshDuration); err != nil {
								logger.ErrorCtx(ctx, "failed to rotate token cookies", zap.Error(err))
								http.Error(w, "unauthorized", http.StatusUnauthorized)
								return
							}
							userID = refreshClaims.UserID
							email = refreshClaims.Email
							found = true
							logger.InfoCtx(ctx, "token rotation completed",
								zap.String("user_id", userID.String()),
							)
						} else {
							logger.WarnCtx(ctx, "refresh token verification failed", zap.Error(err))
						}
					}
				} else {
					logger.WarnCtx(ctx, "access token verification failed", zap.Error(err))
				}
			}

			if !found {
				logger.WarnCtx(ctx, "authentication failed: no valid credentials found",
					zap.String("path", r.URL.Path),
				)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, UserIDKey, userID)
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
