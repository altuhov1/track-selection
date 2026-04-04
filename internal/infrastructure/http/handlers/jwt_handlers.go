package handlers

import (
	"encoding/json"
	"net/http"
	"track-selection/internal/application/auth"
	"track-selection/internal/domain/shared/errors"
)

type AuthHandler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
}

func NewAuthHandler(registerUC *auth.RegisterUseCase, loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Валидация входных данных (быстрая, не бизнес-логика)
	if req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password required")
		return
	}

	if req.Role == "" {
		req.Role = "student"
	}

	// Вызываем Use Case
	err := h.registerUC.Execute(r.Context(), auth.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		// Обработка ошибок
		switch err {
		case errors.ErrAlreadyExists:
			sendError(w, http.StatusConflict, "EMAIL_EXISTS", "email already exists")
		case errors.ErrInvalidEmail:
			sendError(w, http.StatusBadRequest, "INVALID_EMAIL", "invalid email format")
		default:
			if err.Error() == "password must be at least 6 characters" {
				sendError(w, http.StatusBadRequest, "WEAK_PASSWORD", err.Error())
			} else {
				sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
			}
		}
		return
	}

	sendJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "email and password required")
		return
	}

	// Вызываем Use Case
	output, err := h.loginUC.Execute(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"token": output.Token})
}
