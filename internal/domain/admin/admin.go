package admin

import (
	"time"
	"track-selection/internal/domain/shared/value_objects"
)

type Admin struct {
	id         AdminID
	authUserID string
	email      value_objects.Email
	firstName  string
	lastName   string
	createdAt  time.Time
	updatedAt  time.Time
}

func NewAdmin(authUserID string, emailStr string, firstName, lastName string) (*Admin, error) {
	email, err := value_objects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Admin{
		id:         NewAdminID(),
		authUserID: authUserID,
		email:      email,
		firstName:  firstName,
		lastName:   lastName,
		createdAt:  now,
		updatedAt:  now,
	}, nil
}

func NewAdminFromDB(
	id AdminID,
	authUserID string,
	email value_objects.Email,
	firstName, lastName string,
	createdAt time.Time,
	updatedAt time.Time,
) *Admin {
	return &Admin{
		id:         id,
		authUserID: authUserID,
		email:      email,
		firstName:  firstName,
		lastName:   lastName,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
	}
}

func (a *Admin) ID() AdminID                { return a.id }
func (a *Admin) AuthUserID() string         { return a.authUserID }
func (a *Admin) Email() value_objects.Email { return a.email }
func (a *Admin) FirstName() string          { return a.firstName }
func (a *Admin) LastName() string           { return a.lastName }
func (a *Admin) CreatedAt() time.Time       { return a.createdAt }
func (a *Admin) UpdatedAt() time.Time       { return a.updatedAt }
