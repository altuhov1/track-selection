package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/auth"
	"track-selection/internal/application/student"
	"track-selection/internal/application/track"
)

type Handler struct {
	registerUC             *auth.RegisterUseCase
	loginUC                *auth.LoginUseCase
	updatePreferencesUC    *student.UpdatePreferencesUseCase
	getPreferencesUC       *student.GetPreferencesUseCase
	getProfileCompletionUC *student.GetProfileCompletionUseCase
	getAllTracksUC         *track.GetAllTracksUseCase
	createTrackUC          *track.CreateTrackUseCase
	updateTrackUC          *track.UpdateTrackUseCase
	deleteTrackUC          *track.DeleteTrackUseCase
	getRecommendationsUC   *student.GetRecommendationsUseCase
	selectTrackUC          *student.SelectTrackUseCase
	getSelectedTracksUC    *student.GetSelectedTracksUseCase
	unselectTrackUC        *student.UnselectTrackUseCase
}

func NewHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	updatePreferencesUC *student.UpdatePreferencesUseCase,
	getPreferencesUC *student.GetPreferencesUseCase,
	getProfileCompletionUC *student.GetProfileCompletionUseCase,
	getAllUC *track.GetAllTracksUseCase,
	createUC *track.CreateTrackUseCase,
	updateUC *track.UpdateTrackUseCase,
	deleteUC *track.DeleteTrackUseCase,
	getRecommendationsUC *student.GetRecommendationsUseCase,
	selectTrackUC *student.SelectTrackUseCase,
	getSelectedTracksUC *student.GetSelectedTracksUseCase,
	unselectTrackUC *student.UnselectTrackUseCase,
) *Handler {
	return &Handler{
		registerUC:             registerUC,
		loginUC:                loginUC,
		updatePreferencesUC:    updatePreferencesUC,
		getPreferencesUC:       getPreferencesUC,
		getProfileCompletionUC: getProfileCompletionUC,
		getAllTracksUC:         getAllUC,
		createTrackUC:          createUC,
		updateTrackUC:          updateUC,
		deleteTrackUC:          deleteUC,
		getRecommendationsUC:   getRecommendationsUC,
		selectTrackUC:          selectTrackUC,
		getSelectedTracksUC:    getSelectedTracksUC,
		unselectTrackUC:        unselectTrackUC,
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
