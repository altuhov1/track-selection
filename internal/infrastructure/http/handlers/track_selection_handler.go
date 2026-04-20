package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/student"
	"track-selection/internal/domain/shared/errors"

	"github.com/gorilla/mux"
)

func (h *Handler) SelectTrack(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	var input student.SelectTrackInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if input.TrackID == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "track_id is required")
		return
	}

	if err := h.selectTrackUC.Execute(r.Context(), userID, input); err != nil {
		switch err {
		case errors.ErrNotFound:
			sendError(w, http.StatusNotFound, "NOT_FOUND", "student or track not found")
		default:
			sendError(w, http.StatusInternalServerError, "SELECT_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "track selected"})
}

func (h *Handler) GetSelectedTracks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	output, err := h.getSelectedTracksUC.Execute(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		return
	}

	sendJSON(w, http.StatusOK, output)
}

func (h *Handler) UnselectTrack(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "user not found")
		return
	}

	vars := mux.Vars(r)
	trackID := vars["id"]

	if trackID == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "track_id is required")
		return
	}

	if err := h.unselectTrackUC.Execute(r.Context(), userID, trackID); err != nil {
		switch err {
		case errors.ErrNotFound:
			sendError(w, http.StatusNotFound, "NOT_FOUND", "selection not found")
		default:
			sendError(w, http.StatusInternalServerError, "UNSELECT_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "track unselected"})
}
