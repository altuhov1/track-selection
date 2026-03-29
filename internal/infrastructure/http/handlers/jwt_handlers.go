package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	if req.Email == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}
	if req.Password == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}
	if req.Role == "" {
		req.Role = "user"
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		switch err.Error() {
		case "invalid email format":
			sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		case "password must be at least 6 characters":
			sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		case "role must be admin or user":
			sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		case "email already exists":
			sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		default:
			sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	resp := models.RegisterResponse{
		User: *user,
	}
	sendJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid email or password")
			return
		}
		sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	resp := models.LoginResponse{
		Token: token,
	}
	sendJSON(w, http.StatusOK, resp)
}
