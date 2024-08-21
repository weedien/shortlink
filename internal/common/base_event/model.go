package base_event

import (
	"time"
)

type AppEvent interface {
	Name() string
	OccurredAt() time.Time
}

type CommonEvent struct {
	occurredAt time.Time
}

func NewCommonEvent() CommonEvent {
	return CommonEvent{
		occurredAt: time.Now(),
	}
}

func (e CommonEvent) OccurredAt() time.Time {
	return e.occurredAt
}
