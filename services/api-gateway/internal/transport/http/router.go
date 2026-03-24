package router

import (
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/handlers"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
)

// NewRouter wires the chi router with all middlewares and routes.
func NewRouter(cfg *config.Config, db *pgxpool.Pool, logger *zap.Logger, authSvc service.AuthService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(chimiddleware.Recoverer)

	health := &handlers.HealthHandler{DB: db}
	auth := &handlers.AuthHandler{Service: authSvc}

	r.Get("/healthz", health.Healthz)
	r.Get("/readyz", health.Readyz)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Post("/auth/register", auth.Register)
		v1.Post("/auth/login", auth.Login)
		v1.Post("/auth/logout", auth.Logout)

		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(cfg.JWTSecret))
			protected.Get("/me", auth.Me)
		})
	})

	return r
}
