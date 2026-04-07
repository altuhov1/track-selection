package eventbus

import (
	"context"
	"sync"
	"time"
	"track-selection/internal/domain/shared/events"
)

type MemoryEventBus struct {
	handlers map[string][]events.EventHandler
	mu       sync.RWMutex
}

func NewMemoryBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]events.EventHandler),
	}
}

func (bus *MemoryEventBus) Subscribe(eventType string, handler events.EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	return nil
}

func (bus *MemoryEventBus) Publish(ctx context.Context, event events.DomainEvent) error {
	bus.mu.RLock()
	handlers := bus.handlers[event.GetEventType()]
	bus.mu.RUnlock()

	for _, handler := range handlers {
		go func(h events.EventHandler) {

			bgCtx := context.WithoutCancel(ctx)

			// Добавляем свой таймаут для безопасности
			timeoutCtx, cancel := context.WithTimeout(bgCtx, 30*time.Second)
			defer cancel()

			h.Handle(timeoutCtx, event)
		}(handler)
	}

	return nil
}
