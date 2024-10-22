package event

import (
	"shortlink/internal/common/base_event"
	"shortlink/internal/link/domain/link"
)

type UserVisitEvent struct {
	base_event.CommonEvent
	visitInfo link.UserVisitInfo
}

func (e UserVisitEvent) Name() string {
	return "UserVisitEvent"
}

func NewUserVisitEvent(recordInfo link.UserVisitInfo) UserVisitEvent {
	return UserVisitEvent{
		CommonEvent: base_event.NewCommonEvent(),
		visitInfo:   recordInfo,
	}
}
