package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"track-selection/internal/application/auth"
	"track-selection/internal/domain/shared/errors"
)

type AuthHandler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
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
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		fmt.Println(err)
		// Обработка ошибок
		switch err {
		case errors.ErrAlreadyExists:
			sendError(w, http.StatusConflict, "EMAIL_EXISTS", "email already exists")
		case errors.ErrInvalidEmail:
			sendError(w, http.StatusBadRequest, "INVALID_EMAIL", "invalid email format")
		case errors.ErrInvalidRole:
			sendError(w, http.StatusBadRequest, "INVALID_ROLE", "invalid role")
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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(err)
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"token": output.Token})
}

// GetMe возвращает информацию о текущем пользователе
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Данные уже в контексте от middleware!
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "not authenticated")
		return
	}

	role, ok := r.Context().Value("user_role").(string)
	if !ok || role == "" {
		sendError(w, http.StatusUnauthorized, "UNAUTHORIZED", "role not found in context")
		return
	}

	firstName, _ := r.Context().Value("first_name").(string)
	lastName, _ := r.Context().Value("last_name").(string)
	email, _ := r.Context().Value("email").(string)

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"id":         userID,
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
	})
}
