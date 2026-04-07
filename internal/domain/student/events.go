package student

import (
	"time"
	"track-selection/internal/domain/shared/events"

	"github.com/google/uuid"
)

type StudentRegisteredEvent struct {
	events.BaseDomainEvent
	UserID string
	Email  string
}

// NewUserRegisteredEvent создает событие регистрации пользователя
func NewStudentRegisteredEvent(userID, email, role string) StudentRegisteredEvent {
	return StudentRegisteredEvent{
		BaseDomainEvent: events.BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  "student.registered",
			OccurredAt: time.Now(),
		},
		UserID: userID,
		Email:  email,
	}
}
