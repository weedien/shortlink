package link

import (
	"context"
	"shortlink/common/types"
	"shortlink/internal/domain/link/aggregate"
	"shortlink/internal/domain/link/valobj"
)

type Repository interface {
	ShortUriExists(ctx context.Context, shortUrl string) (bool, error)

	CreateLink(ctx context.Context, aggregate aggregate.CreateLinkAggregate) error

	//CreateLinkWithLock(ctx context.Context, aggregate aggregate.CreateLinkAggregate) error

	UpdateLink(
		ctx context.Context,
		id types.LinkID,
		enableStatus int,
		updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
	) error

	GetOriginalUrlByShortUrl(
		ctx context.Context,
		fullShortUrl string,
		statsInfo valobj.ShortLinkStatsRecordVo,
	) (string, error)

	//RecordLinkVisitInfo(ctx context.Context, info valobj.ShortLinkStatsRecordVo) error
}
