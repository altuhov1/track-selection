// internal/domain/shared/events/user_events.go
package events

import (
	"time"

	"github.com/google/uuid"
)

// UserRegisteredEvent — событие регистрации пользователя
type UserRegisteredEvent struct {
	BaseDomainEvent
	UserID string
	Email  string
	Role   string
}

// NewUserRegisteredEvent создает событие регистрации пользователя
func NewUserRegisteredEvent(userID, email, role string) UserRegisteredEvent {
	return UserRegisteredEvent{
		BaseDomainEvent: BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  "student.registered",
			OccurredAt: time.Now(),
		},
		UserID: userID,
		Email:  email,
		Role:   role,
	}
}
