package server

import (
	"log/slog"
	"net/http"

	"github.com/ParkerGits/go-backend-starter/pkg/httplib"
	"github.com/rs/cors"
)

func (r *router) withCors(next http.Handler) http.Handler {
	logger := slog.Default()
	return cors.New(cors.Options{
		AllowedOrigins:   r.config.CorsAllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Logger:           slog.NewLogLogger(logger.Handler(), slog.LevelDebug),
	}).Handler(next)
}

func (r *router) withDefaultMiddleware(next http.Handler) http.Handler {
	return httplib.CreateStack(
		httplib.WithRequestLogger,
		r.withCors,
	)(next)
}
