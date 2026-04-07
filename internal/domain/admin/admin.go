package admin

import (
	"errors"
	"strings"
	"time"
	"track-selection/internal/domain/shared/value_objects"
)

type Admin struct {
	id         AdminID
	authUserID string
	email      value_objects.Email
	username   string
	rating     int
	createdAt  time.Time
	updatedAt  time.Time
}

func NewAdmin(authUserID string, emailStr string) (*Admin, error) {
	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Admin{
		id:         NewAdminID(),
		authUserID: authUserID,
		email:      email,
		username:   strings.Split(emailStr, "@")[0], // username из email по умолчанию
		rating:     0,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func NewAdminFromDB(
	id AdminID,
	authUserID string,
	email value_objects.Email,
	username string,
	rating int,
	createdAt time.Time,
	updatedAt time.Time,
) *Admin {
	return &Admin{
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
func (s *Admin) ID() AdminID                { return s.id }
func (s *Admin) AuthUserID() string         { return s.authUserID }
func (s *Admin) Email() value_objects.Email { return s.email }
func (s *Admin) Username() string           { return s.username }
func (s *Admin) Rating() int                { return s.rating }
func (s *Admin) CreatedAt() time.Time       { return s.createdAt }
func (s *Admin) UpdatedAt() time.Time       { return s.updatedAt }

// Методы для изменения данных
func (s *Admin) ChangeUsername(username string) error {
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

func (s *Admin) ChangeRating(rating int) error {
	if rating < 0 || rating > 100 {
		return errors.New("rating must be between 0 and 100")
	}

	s.rating = rating
	s.updatedAt = time.Now()
	return nil
}

func (s *Admin) ChangeEmail(email value_objects.Email) {
	s.email = email
	s.updatedAt = time.Now()
}
