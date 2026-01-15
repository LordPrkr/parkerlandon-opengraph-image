package httplib

import (
	"log/slog"
	"net/http"
	"time"
)

// Applies argument middleware from top down (top is innermost, bottom is outermost)
func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

func WithRequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := slog.With(slog.Group("request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		))
		slog.SetDefault(requestLogger)
		requestLogger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"incoming request",
		)

		start := time.Now()
		wrapped := &WrappedWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)

		responseLogger := requestLogger.With(slog.Group("response",
			slog.Int("status", wrapped.StatusCode),
			slog.String("elapsed_time", time.Since(start).String()),
		))
		slog.SetDefault(responseLogger)
		responseLogger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"request complete",
		)
	})
}
