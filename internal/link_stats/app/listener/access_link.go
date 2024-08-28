package listener

import (
	"context"
	"shortlink/internal/common/base_event"
	"shortlink/internal/link/app/event"
	"shortlink/internal/link_stats/domain"
)

type RecordLinkVisitListener struct {
	repo domain.Repository
}

func NewRecordLinkVisitListener(repo domain.Repository) RecordLinkVisitListener {
	return RecordLinkVisitListener{repo: repo}
}

func (h RecordLinkVisitListener) Process(ctx context.Context, e base_event.Event) error {
	if ve, ok := e.(event.RecordLinkVisitEvent); ok {
		return h.repo.SaveLinkStats(ctx, ve.RecordInfo)
	}
	return nil
}
