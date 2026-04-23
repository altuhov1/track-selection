package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"track-selection/internal/domain/auth"
	domErr "track-selection/internal/domain/shared/errors"
)

const (
	RoleAny   int = 0 // любая роль (admin или user)
	RoleAdmin int = 1 // только admin
	RoleUser  int = 2 // только user
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
			switch {
			case errors.Is(err, domErr.ErrTokenExpired):
				writeJSONError(w, http.StatusUnauthorized, "TOKEN_EXPIRED", "token has expired")
			case errors.Is(err, domErr.ErrInvalidToken):
				writeJSONError(w, http.StatusUnauthorized, "INVALID_TOKEN", "invalid token")
			default:
				writeJSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
			}
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
		}
		ctx := context.WithValue(r.Context(), "user_id", token.UserID)
		ctx = context.WithValue(ctx, "user_role", string(token.Role))
		ctx = context.WithValue(ctx, "first_name", token.FirstName)
		ctx = context.WithValue(ctx, "last_name", token.LastName)
		ctx = context.WithValue(ctx, "email", token.Email)
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
