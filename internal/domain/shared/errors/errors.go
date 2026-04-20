package errors

import "errors"

var (
	ErrNotFound             = errors.New("entity not found")
	ErrInvalidEmail         = errors.New("invalid email address")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrAlreadyExists        = errors.New("entity already exists")
	ErrInvalidRole          = errors.New("we have not this role")
	ErrTokenExpired         = errors.New("token has expired")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidRequest       = errors.New("invalid request")
	ErrInvalidGrade         = errors.New("grade must be between 2 and 5")
	ErrInvalidSkillValue    = errors.New("val must be between 1 and 10")
	ErrInvalidLearningStyle = errors.New("val must be between 1 and 3")
	ErrInvalidCertificate   = errors.New("val must be between 0 and 1")
	ErrProfileNotComplete   = errors.New("profile is not complete enough for recommendations")
)
