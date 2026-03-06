package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/http/response"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func (h HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}

func (h HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]any{
			"ok":    false,
			"error": "db_not_ready",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"ok": true,
	})
}