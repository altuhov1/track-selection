package student

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Student struct {
	id        StudentID
	email     Email
	username  string
	createdAt time.Time
	updatedAt time.Time
}

func NewStudent(emailStr string, username string) (*Student, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username is required")
	}

	if len(username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	now := time.Now()

	return &Student{
		id:        NewStudentID(),
		email:     email,
		username:  username,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (s *Student) ID() StudentID {
	return s.id
}

func (s *Student) Email() Email {
	return s.email
}

func (s *Student) Username() string {
	return s.username
}
