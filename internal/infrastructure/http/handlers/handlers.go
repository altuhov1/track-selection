package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/auth"
)

type Handler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
}

func NewHandler(registerUC *auth.RegisterUseCase, loginUC *auth.LoginUseCase) *Handler {
	return &Handler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
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
