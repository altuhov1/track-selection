package events

import "context"

type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
}
