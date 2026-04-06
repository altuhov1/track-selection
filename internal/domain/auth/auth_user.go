package auth

import (
	"errors"
	"time"
	"track-selection/internal/domain/shared/value_objects"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleAdmin   UserRole = "admin"
)

type AuthUser struct {
	ID           string
	Email        value_objects.Email
	PasswordHash string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAuthUser(email string, rawPassword string, role UserRole) (*AuthUser, error) {
	validatedEmail, err := value_objects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	if len(rawPassword) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &AuthUser{
		ID:           uuid.New().String(),
		Email:        validatedEmail,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *AuthUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
