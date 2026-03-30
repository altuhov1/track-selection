package student

import (
	"errors"
	"net/mail"
	"strings"
	"track-selection/internal/domain/shared"

	"github.com/google/uuid"
)

// ---------- Email ----------
type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return Email{}, shared.ErrInvalidEmail
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, shared.ErrInvalidEmail
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return Email{}, shared.ErrInvalidEmail
	}

	if len(parts[0]) > 64 || len(parts[1]) > 255 {
		return Email{}, shared.ErrInvalidEmail
	}

	return Email{value: addr.Address}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return strings.EqualFold(e.value, other.value)
}

// ---------- StudentID ----------
type StudentID struct {
	value string
}

func NewStudentID() StudentID {
	return StudentID{value: uuid.New().String()}
}

func StudentIDFromString(id string) (StudentID, error) {
	if id == "" {
		return StudentID{}, errors.New("student id cannot be empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		return StudentID{}, errors.New("invalid student id format")
	}

	return StudentID{value: id}, nil
}

func (id StudentID) String() string {
	return id.value
}

func (id StudentID) Equals(other StudentID) bool {
	return id.value == other.value
}

func (id StudentID) IsEmpty() bool {
	return id.value == ""
}
