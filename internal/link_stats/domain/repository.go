package domain

import (
	"context"
	"shortlink/internal/link_stats/domain/event"
)

type Repository interface {
	SaveLinkStat(ctx context.Context, statsInfo event.UserVisitInfo) error
}
