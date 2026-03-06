package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "request_id"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			// Simple request id; can upgrade later (uuid)
			reqID = time.Now().UTC().Format("20060102T150405.000000000Z07:00")
		}

		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), CtxKeyRequestID, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) string {
	v := ctx.Value(CtxKeyRequestID)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func WithLogger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status
			rw := &statusWriter{ResponseWriter: w, status: 200}
			next.ServeHTTP(rw, r)

			reqID := RequestIDFromContext(r.Context())

			logger.Info().
				Str("request_id", reqID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", rw.status).
				Dur("duration_ms", time.Since(start)).
				Msg("request")
		})
	}
}

func Recoverer(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					reqID := RequestIDFromContext(r.Context())
					logger.Error().Str("request_id", reqID).Any("panic", rec).Msg("panic recovered")
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}