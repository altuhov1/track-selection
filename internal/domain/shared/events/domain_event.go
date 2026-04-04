package events

import "time"

type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetOccurredAt() time.Time
}

type BaseDomainEvent struct {
	EventID    string
	EventType  string
	OccurredAt time.Time
}

func (e BaseDomainEvent) GetEventID() string       { return e.EventID }
func (e BaseDomainEvent) GetEventType() string     { return e.EventType }
func (e BaseDomainEvent) GetOccurredAt() time.Time { return e.OccurredAt }
