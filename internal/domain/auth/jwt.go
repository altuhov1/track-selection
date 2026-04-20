package auth

type JWTService interface {
	GenerateToken(userID string, role UserRole, firstName, lastName, email string) (string, error)

	ValidateToken(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID    string
	Role      UserRole
	FirstName string
	LastName  string
	Email     string
}
