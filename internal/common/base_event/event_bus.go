package base_event

import (
	"context"
	"reflect"
)

// EventBus is an interface for event bus in the application layer.
type EventBus interface {
	Publish(ctx context.Context, event Event)
	Subscribe(eventType reflect.Type, consumer EventListener)
}
