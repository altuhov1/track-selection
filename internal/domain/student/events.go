package student

import (
	"time"
	"track-selection/internal/domain/shared/events"

	"github.com/google/uuid"
)

type StudentRegisteredEvent struct {
	events.BaseDomainEvent
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

// NewStudentRegisteredEvent создает событие регистрации студента
func NewStudentRegisteredEvent(userID, email, firstName, lastName string) StudentRegisteredEvent {
	return StudentRegisteredEvent{
		BaseDomainEvent: events.BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  "student.registered",
			OccurredAt: time.Now(),
		},
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}
}

// Это пометка что пользователь готов к показу статистик
type ProfileCompletedEvent struct {
	events.BaseDomainEvent
	UserID    string `json:"user_id"`
	Completed bool   `json:"completed"`
}

func NewProfileCompletedEvent(userID string, completed bool) ProfileCompletedEvent {
	eventType := "profile.completed"
	if !completed {
		eventType = "profile.incomplete"
	}

	return ProfileCompletedEvent{
		BaseDomainEvent: events.BaseDomainEvent{
			EventID:    uuid.New().String(),
			EventType:  eventType,
			OccurredAt: time.Now(),
		},
		UserID:    userID,
		Completed: completed,
	}
}
