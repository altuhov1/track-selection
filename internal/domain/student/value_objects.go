package student

import (
	"errors"

	"github.com/google/uuid"
)

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
