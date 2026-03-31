package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/transport/respond"
)

type ctxKey int

const (
	ctxKeyUserID   ctxKey = iota
	ctxKeySchoolID ctxKey = iota
)

// RequireAuth validates JWT access tokens issued by api-gateway,
// extracts user_id and school_id claims, and injects them into context.
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			tokenType, _ := claims["type"].(string)
			if tokenType != "access" {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userIDStr, _ := claims["user_id"].(string)
			schoolIDStr, _ := claims["school_id"].(string)

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			schoolID, err := uuid.Parse(schoolIDStr)
			if err != nil {
				respond.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
			ctx = context.WithValue(ctx, ctxKeySchoolID, schoolID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromCtx extracts the authenticated user's ID from context.
func UserIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKeyUserID).(uuid.UUID)
	return id, ok
}

// SchoolIDFromCtx extracts the authenticated user's school ID from context.
func SchoolIDFromCtx(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKeySchoolID).(uuid.UUID)
	return id, ok
}
