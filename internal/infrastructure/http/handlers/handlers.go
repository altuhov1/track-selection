package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/auth"
	"track-selection/internal/application/student"
)

type Handler struct {
	registerUC             *auth.RegisterUseCase
	loginUC                *auth.LoginUseCase
	updatePreferencesUC    *student.UpdatePreferencesUseCase
	getPreferencesUC       *student.GetPreferencesUseCase
	getProfileCompletionUC *student.GetProfileCompletionUseCase
}

func NewHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	updatePreferencesUC *student.UpdatePreferencesUseCase,
	getPreferencesUC *student.GetPreferencesUseCase,
	getProfileCompletionUC *student.GetProfileCompletionUseCase, // ← добавить
) *Handler {
	return &Handler{
		registerUC:             registerUC,
		loginUC:                loginUC,
		updatePreferencesUC:    updatePreferencesUC,
		getPreferencesUC:       getPreferencesUC,
		getProfileCompletionUC: getProfileCompletionUC, // ← добавить
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
