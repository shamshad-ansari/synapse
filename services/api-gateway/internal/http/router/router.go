package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/handlers"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/middleware"
)

func New(cfg config.Config, db *pgxpool.Pool, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()

	// Middleware order matters
	r.Use(middleware.WithRequestID)
	r.Use(middleware.Recoverer(logger))
	r.Use(middleware.WithLogger(logger))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.CORSAllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	health := handlers.HealthHandler{DB: db}
	me := handlers.MeHandler{}

	r.Get("/healthz", health.Healthz)
	r.Get("/readyz", health.Readyz)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/me", me.GetMe)
	})

	return r
}