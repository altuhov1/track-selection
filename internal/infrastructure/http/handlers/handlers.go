package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	jwtService services.JWTservice
}

func NewHandler(
	jwtService services.JWTservice,

) (*Handler, error) {
	return &Handler{
		jwtService: jwtService,
	}, nil
}

func (h *Handler) TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	slog.Info("Use test handler")
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
