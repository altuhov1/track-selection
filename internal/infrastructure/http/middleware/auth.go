package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"track-selection/internal/domain/shared/auth"
)

const (
	UserIDKey string = "user_id"
	RoleKey   string = "role"
	RoleAny   int    = 0 // любая роль (admin или user)
	RoleAdmin int    = 1 // только admin
	RoleUser  int    = 2 // только user
)

func WithAuth(JwtService auth.JWTService, handlerFn func(http.ResponseWriter, *http.Request), requiredRole int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
			return
		}

		tokenString := parts[1]
		token, err := JwtService.ValidateToken(tokenString)

		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			return
		}

		if token.UserID == "" {
			writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid token")
			return
		}

		switch requiredRole {
		case RoleAdmin:
			if token.Role != "admin" {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "access denied: admin role required")
				return
			}
		case RoleUser:
			if token.Role != "user" {
				writeJSONError(w, http.StatusForbidden, "FORBIDDEN", "access denied: user role required")
				return
			}
		case RoleAny:
			// любая роль подходит
		}
		//nolint:staticcheck
		ctx := context.WithValue(r.Context(), UserIDKey, token.UserID)
		//nolint:staticcheck
		ctx = context.WithValue(ctx, RoleKey, token.Role)
		r = r.WithContext(ctx)

		handlerFn(w, r)
	}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
