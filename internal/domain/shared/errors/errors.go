package errors

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidEmail  = errors.New("invalid email address")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrInvalidRole   = errors.New("we have not this role")
	ErrTokenExpired  = errors.New("token has expired")
	ErrInvalidToken  = errors.New("invalid token")
)
