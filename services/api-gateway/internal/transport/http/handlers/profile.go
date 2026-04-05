package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// ProfileHandler serves GET /v1/profile/summary.
type ProfileHandler struct {
	Service service.ProfileService
}

// Summary returns aggregated profile metrics for the current user.
func (h *ProfileHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	schoolID, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	out, err := h.Service.GetSummary(r.Context(), userID, schoolID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	respond.JSON(w, http.StatusOK, out)
}
