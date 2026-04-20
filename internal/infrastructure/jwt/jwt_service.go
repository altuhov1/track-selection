package jwt

import (
	"fmt"
	"time"

	"errors"
	"track-selection/internal/domain/auth"
	domErr "track-selection/internal/domain/shared/errors"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type JWTServiceImpl struct {
	secret     string
	expiration time.Duration
}

func NewJWTService(config *JWTConfig) auth.JWTService {
	return &JWTServiceImpl{
		secret:     config.Secret,
		expiration: config.Expiration,
	}
}

func (s *JWTServiceImpl) GenerateToken(userID string, role auth.UserRole, firstName, lastName, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"role":       string(role),
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"exp":        time.Now().Add(s.expiration).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domErr.ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, domErr.ErrInvalidToken
		}
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, domErr.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return nil, domErr.ErrInvalidToken
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return nil, domErr.ErrInvalidToken
	}
	first_name, ok := claims["first_name"].(string)
	if !ok || userID == "" {
		return nil, domErr.ErrInvalidToken
	}

	last_name, ok := claims["last_name"].(string)
	if !ok {
		return nil, domErr.ErrInvalidToken
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, domErr.ErrInvalidToken
	}

	return &auth.TokenClaims{
		UserID:    userID,
		Role:      auth.UserRole(roleStr),
		FirstName: first_name,
		LastName:  last_name,
		Email:     email,
	}, nil
}
