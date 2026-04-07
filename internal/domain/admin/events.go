package admin

import (
	"time"
	"track-selection/internal/domain/shared/events"

	"github.com/google/uuid"
)

type AdminRegisteredEvent struct {
	events.BaseDomainEvent
	UserID string
	Email  string
}

func NewAdminRegisteredEvent(userID, email, role string) AdminRegisteredEvent {
	return AdminRegisteredEvent{
		BaseDomainEvent: events.BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  "admin.registered",
			OccurredAt: time.Now(),
		},
		UserID: userID,
		Email:  email,
	}
}
