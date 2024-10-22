package base_event

import (
	"github.com/google/uuid"
	"time"
)

type Event interface {
	Id() string
	Name() string
	OccurredAt() time.Time
}

type CommonEvent struct {
	id         string
	occurredAt time.Time
}

func NewCommonEvent() CommonEvent {
	return CommonEvent{
		id:         uuid.New().String(),
		occurredAt: time.Now(),
	}
}

func (e CommonEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e CommonEvent) Id() string {
	return e.id
}

func (e CommonEvent) Name() string {
	return "CommonEvent"
}
