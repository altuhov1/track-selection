package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/domain/shared/errors"
)

func (h *Handler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	var updates json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.updatePreferencesUC.Execute(r.Context(), userID, updates); err != nil {
		// Обрабатываем различные типы ошибок
		switch err {
		case errors.ErrInvalidRequest:
			sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		case errors.ErrInvalidLearningStyle:
			sendError(w, http.StatusBadRequest, "INVALID_LEARNING_STYLE", err.Error())
		case errors.ErrInvalidCertificate:
			sendError(w, http.StatusBadRequest, "INVALID_CERTIFICATE", err.Error())
		case errors.ErrInvalidSkillValue:
			sendError(w, http.StatusBadRequest, "INVALID_SKILL_VALUE", err.Error())
		case errors.ErrInvalidGrade:
			sendError(w, http.StatusBadRequest, "INVALID_GRADE", err.Error())
		default:
			sendError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "preferences updated"})
}

func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	prefs, err := h.getPreferencesUC.Execute(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get preferences")
		return
	}

	sendJSON(w, http.StatusOK, prefs)
}

func (h *Handler) GetProfileCompletion(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	status, err := h.getProfileCompletionUC.Execute(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to get profile status")
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"is_complete":  status.IsComplete,
		"completed_at": status.CompletedAt,
	})
}
