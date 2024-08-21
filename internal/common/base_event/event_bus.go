package base_event

import (
	"context"
	"log/slog"
	"reflect"
	"sync"
)

// AppEventBus is an interface for event bus in the application layer.
type AppEventBus interface {
	Publish(ctx context.Context, event AppEvent)
	Subscribe(eventType reflect.Type, consumer AppEventListener)
}

type SimpleEventBus struct {
	consumers map[reflect.Type][]AppEventListener
	mu        sync.RWMutex
}

func NewSimpleEventBus() AppEventBus {
	return &SimpleEventBus{
		consumers: make(map[reflect.Type][]AppEventListener),
	}
}

func (bus *SimpleEventBus) Publish(ctx context.Context, event AppEvent) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	eventType := reflect.TypeOf(event)
	if consumers, found := bus.consumers[eventType]; found {
		for _, consumer := range consumers {
			go func(consumer AppEventListener) {
				if err := consumer.Process(ctx, event); err != nil {
					slog.Info("event consumer error", err)
				}
			}(consumer)
		}
		return
	}
	slog.Warn("no consumer found for event", "event", event.Name())
}

func (bus *SimpleEventBus) Subscribe(eventType reflect.Type, consumer AppEventListener) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.consumers[eventType] = append(bus.consumers[eventType], consumer)
}
