package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/ticketing-service/pkg/logger"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// RequestIDMiddleware adds a unique request ID to the context
// This comes before LoggingMiddleware so the ID is available for all logs
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		// Add to context for logger
		ctx := logger.WithRequestID(r.Context(), requestID)

		// Also add to response header for client-side debugging
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs HTTP requests and responses with structured logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Wrap response writer to capture status code and size
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Log incoming request
		logger.InfoCtx(ctx, "incoming request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
		)

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log response with duration
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.status),
			zap.Int("size_bytes", wrapped.size),
			zap.Duration("duration", duration),
			zap.Float64("duration_ms", float64(duration.Nanoseconds())/1e6),
		}

		// Use different log levels based on status code
		switch {
		case wrapped.status >= 500:
			logger.ErrorCtx(ctx, "request completed with server error", fields...)
		case wrapped.status >= 400:
			logger.WarnCtx(ctx, "request completed with client error", fields...)
		default:
			logger.InfoCtx(ctx, "request completed", fields...)
		}
	})
}

// RecovererMiddleware wraps chi's Recoverer but adds structured logging
func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				ctx := r.Context()

				logger.ErrorCtx(ctx, "panic recovered",
					zap.Any("panic", rvr),
					zap.String("stack", string(debug.Stack())),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
				)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RateLimitLogger adds logging when rate limits are hit
func RateLimitLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if rate limit header is set (httprate sets this)
		next.ServeHTTP(w, r)

		if w.Header().Get("X-RateLimit-Remaining") == "0" {
			ctx := r.Context()
			logger.WarnCtx(ctx, "rate limit threshold reached",
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("path", r.URL.Path),
			)
		}
	})
}
