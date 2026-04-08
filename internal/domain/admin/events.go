package admin

import (
	"time"
	"track-selection/internal/domain/shared/events"

	"github.com/google/uuid"
)

type AdminRegisteredEvent struct {
	events.BaseDomainEvent
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

func NewAdminRegisteredEvent(userID, email, firstName, lastName string) AdminRegisteredEvent {
	return AdminRegisteredEvent{
		BaseDomainEvent: events.BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  "admin.registered",
			OccurredAt: time.Now(),
		},
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
}
