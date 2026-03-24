package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/service"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

type AuthHandler struct {
	Service service.AuthService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := h.Service.Register(r.Context(), in)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, out)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	out, err := h.Service.Login(r.Context(), in)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, out)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	respond.JSON(w, http.StatusOK, "logged out")
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.Service.GetCurrentUser(r.Context(), userID, schoolID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respond.Error(w, http.StatusNotFound, "user not found")
			return
		}
		respond.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	respond.JSON(w, http.StatusOK, user.ToResponse())
}

func mapAuthError(w http.ResponseWriter, err error) {
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		respond.Error(w, http.StatusUnprocessableEntity, validationErr.Message)
		return
	}
	if errors.Is(err, domain.ErrConflict) {
		respond.Error(w, http.StatusConflict, "email already registered at this school")
		return
	}
	if errors.Is(err, domain.ErrInvalidCredentials) {
		respond.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	respond.Error(w, http.StatusInternalServerError, "internal error")
}
