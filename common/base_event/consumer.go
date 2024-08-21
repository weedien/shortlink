package base_event

import (
	"context"
)

// AppEventListener is an interface for consuming events in the application layer.
type AppEventListener interface {
	Process(ctx context.Context, e AppEvent) error
}
