package student

import (
	"errors"
	"strings"
	"time"
	"track-selection/internal/domain/shared/value_objects"
)

type Student struct {
	id         StudentID
	authUserID string
	email      value_objects.Email
	username   string
	rating     int
	createdAt  time.Time
	updatedAt  time.Time
}

// NewStudent создает нового студента
func NewStudent(authUserID string, emailStr string) (*Student, error) {
	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Student{
		id:         NewStudentID(),
		authUserID: authUserID,
		email:      email,
		username:   strings.Split(emailStr, "@")[0], // username из email по умолчанию
		rating:     0,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

// NewStudentFromDB восстанавливает студента из БД
func NewStudentFromDB(
	id StudentID,
	authUserID string,
	email value_objects.Email,
	username string,
	rating int,
	createdAt time.Time,
	updatedAt time.Time,
) *Student {
	return &Student{
		id:         id,
		authUserID: authUserID,
		email:      email,
		username:   username,
		rating:     rating,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

// Геттеры
func (s *Student) ID() StudentID              { return s.id }
func (s *Student) AuthUserID() string         { return s.authUserID }
func (s *Student) Email() value_objects.Email { return s.email }
func (s *Student) Username() string           { return s.username }
func (s *Student) Rating() int                { return s.rating }
func (s *Student) CreatedAt() time.Time       { return s.createdAt }
func (s *Student) UpdatedAt() time.Time       { return s.updatedAt }

// Методы для изменения данных
func (s *Student) ChangeUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	s.username = username
	s.updatedAt = time.Now()
	return nil
}

func (s *Student) ChangeRating(rating int) error {
	if rating < 0 || rating > 100 {
		return errors.New("rating must be between 0 and 100")
	}

	s.rating = rating
	s.updatedAt = time.Now()
	return nil
}

func (s *Student) ChangeEmail(email value_objects.Email) {
	s.email = email
	s.updatedAt = time.Now()
}
