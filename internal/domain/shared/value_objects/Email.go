package value_objects

import (
	"net/mail"
	"strings"
	"track-selection/internal/domain/shared/errors"
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return Email{}, errors.ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, errors.ErrInvalidEmail
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return Email{}, errors.ErrInvalidEmail
	}

	if len(parts[0]) > 64 || len(parts[1]) > 255 {
		return Email{}, errors.ErrInvalidEmail
	}

	return Email{value: addr.Address}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return strings.EqualFold(e.value, other.value)
}
