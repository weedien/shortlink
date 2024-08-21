package domain

import (
	"context"
	"shortlink/internal/link/domain/valobj"
)

type Repository interface {
	SaveLinkStats(ctx context.Context, statsInfo valobj.ShortLinkStatsRecordVo) error
}
