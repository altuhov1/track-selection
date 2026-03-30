package shared

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidEmail  = errors.New("invalid email address")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrAlreadyExists = errors.New("entity already exists")
)
