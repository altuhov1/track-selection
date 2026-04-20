package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/track"
	"track-selection/internal/domain/shared/errors"

	"github.com/gorilla/mux"
)


func (h *Handler) GetAllTracks(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.getAllTracksUC.Execute(r.Context())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	sendJSON(w, http.StatusOK, tracks)
}

func (h *Handler) CreateTrack(w http.ResponseWriter, r *http.Request) {
	var input track.CreateTrackInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	t, err := h.createTrackUC.Execute(r.Context(), input)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, t)
}

func (h *Handler) UpdateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updates json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if err := h.updateTrackUC.Execute(r.Context(), id, updates); err != nil {
		switch err {
		case errors.ErrNotFound:
			sendError(w, http.StatusNotFound, "NOT_FOUND", "track not found")
		default:
			sendError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "track updated"})
}

func (h *Handler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.deleteTrackUC.Execute(r.Context(), id); err != nil {
		switch err {
		case errors.ErrNotFound:
			sendError(w, http.StatusNotFound, "NOT_FOUND", "track not found")
		default:
			sendError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		}
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"message": "track deleted"})
}
