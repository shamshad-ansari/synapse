package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shamshad-ansari/synapse/api-gateway/internal/transport/respond"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

// Healthz always returns 200. No DB check.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, "ok")
}

// Readyz pings the DB with a 2s timeout.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		respond.Error(w, http.StatusServiceUnavailable, "db unavailable")
		return
	}

	respond.JSON(w, http.StatusOK, "ok")
}
