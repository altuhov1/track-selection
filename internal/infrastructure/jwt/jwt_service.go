package jwt

import (
	"time"

	"track-selection/internal/domain/auth"
	"track-selection/internal/domain/shared/errors"

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

func NewJWTService(config JWTConfig) auth.JWTService {
	return &JWTServiceImpl{
		secret:     config.Secret,
		expiration: config.Expiration,
	}
}

func (s *JWTServiceImpl) GenerateToken(userID string, role auth.UserRole) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    string(role),
		"exp":     time.Now().Add(s.expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrUnauthorized
	}

	return &auth.TokenClaims{
		UserID: claims["user_id"].(string),
		Role:   auth.UserRole(claims["role"].(string)),
	}, nil
}
