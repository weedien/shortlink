package base_event

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
)

type SimpleEventBus struct {
	consumers map[reflect.Type][]EventListener
	mu        sync.RWMutex
}

func NewSimpleEventBus() EventBus {
	return &SimpleEventBus{
		consumers: make(map[reflect.Type][]EventListener),
	}
}

func (bus *SimpleEventBus) Publish(ctx context.Context, event Event) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	eventType := reflect.TypeOf(event)
	if consumers, found := bus.consumers[eventType]; found {
		for _, consumer := range consumers {
			go func(consumer EventListener) {
				if err := consumer.Process(ctx, event); err != nil {
					slog.Info("event consumer error", err)
				}
			}(consumer)
		}
		return
	}
	slog.Warn("no consumer found for event", "event", event.Name())
}

func (bus *SimpleEventBus) Subscribe(eventType reflect.Type, consumer EventListener) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.consumers[eventType] = append(bus.consumers[eventType], consumer)
}
