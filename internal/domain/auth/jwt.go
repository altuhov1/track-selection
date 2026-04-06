package auth

type JWTService interface {
	// GenerateToken создает JWT токен для пользователя
	GenerateToken(userID string, role UserRole) (string, error)

	// ValidateToken проверяет токен и возвращает данные пользователя
	ValidateToken(token string) (*TokenClaims, error)
}

// TokenClaims — данные, которые хранятся в токене
type TokenClaims struct {
	UserID string
	Role   UserRole
}
