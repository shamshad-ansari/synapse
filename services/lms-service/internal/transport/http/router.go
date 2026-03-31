package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/config"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/transport/http/handlers"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/transport/http/middleware"
)

// NewRouter wires the chi router with all middlewares and routes.
func NewRouter(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, logger *zap.Logger, repo domain.LMSRepository) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.Use(chimiddleware.Recoverer)

	health := &handlers.HealthHandler{DB: db}
	canvas := &handlers.CanvasHandler{
		Cfg:    cfg,
		Repo:   repo,
		Redis:  rdb,
		Logger: logger,
	}

	r.Get("/healthz", health.Healthz)
	r.Get("/readyz", health.Readyz)

	r.Route("/v1", func(v1 chi.Router) {
		// OAuth callback — no JWT required (browser redirect from Canvas)
		v1.Get("/lms/callback/canvas", canvas.CallbackCanvas)

		v1.Group(func(protected chi.Router) {
			protected.Use(middleware.RequireAuth(cfg.JWTSecret))

			protected.Get("/lms/connect/canvas", canvas.ConnectCanvas)
			protected.Post("/lms/connect/token", canvas.ConnectToken)
			protected.Get("/lms/status", canvas.Status)
			protected.Post("/lms/sync", canvas.Sync)
			protected.Delete("/lms/disconnect", canvas.Disconnect)
		})
	})

	return r
}
