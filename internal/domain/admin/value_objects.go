package admin

import (
	"errors"

	"github.com/google/uuid"
)

// ---------- AdminID ----------
type AdminID struct {
	value string
}

func NewAdminID() AdminID {
	return AdminID{value: uuid.New().String()}
}

func AdminIDFromString(id string) (AdminID, error) {
	if id == "" {
		return AdminID{}, errors.New("admin id cannot be empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		return AdminID{}, errors.New("invalid admin id format")
	}

	return AdminID{value: id}, nil
}

func (id AdminID) String() string {
	return id.value
}

func (id AdminID) Equals(other AdminID) bool {
	return id.value == other.value
}

func (id AdminID) IsEmpty() bool {
	return id.value == ""
}
