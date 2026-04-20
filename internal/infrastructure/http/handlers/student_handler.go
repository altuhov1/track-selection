package handlers

import (
	"net/http"
	"track-selection/internal/domain/shared/errors"
)

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	output, err := h.getRecommendationsUC.Execute(r.Context(), userID)
	if err != nil {
		switch err {
		case errors.ErrProfileNotComplete:
			sendError(w, http.StatusBadRequest, "PROFILE_NOT_COMPLETE",
				"Please complete your profile (grades and learning style) before getting recommendations")
		default:
			sendError(w, http.StatusInternalServerError, "RECOMMENDATION_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, output)
}
