package listener

import (
	"context"
	"shortlink/common/base_event"
	"shortlink/internal/app/link/event"
	"shortlink/internal/app/link/service"
)

type RecordLinkVisitListener struct {
	service service.LinkAppService
}

func NewLinkAccessedEventListener(service service.LinkAppService) RecordLinkVisitListener {
	if service == nil {
		panic("nil service")
	}

	return RecordLinkVisitListener{service: service}
}

func (h RecordLinkVisitListener) Process(ctx context.Context, e base_event.AppEvent) error {
	if ve, ok := e.(event.RecordLinkVisitEvent); ok {
		return h.service.RecordLinkVisitInfo(ctx, ve.RecordInfo)
	}
	return nil
}
