package domain

import (
	"context"
	"shortlink/internal/link/domain/aggregate"
	"shortlink/internal/link/domain/entity"
	"shortlink/internal/link/domain/valobj"
)

type LinkRepository interface {
	ShortUriExists(ctx context.Context, shortUrl string) (bool, error)

	CreateLink(ctx context.Context, aggregate aggregate.CreateLinkAggregate) error

	//CreateLinkWithLock(ctx context.Context, aggregate aggregate.CreateLinkAggregate) error

	UpdateLink(
		ctx context.Context,
		id entity.LinkID,
		enableStatus int,
		updateFn func(ctx context.Context, link *entity.Link) (*entity.Link, error),
	) error

	GetOriginalUrlByShortUrl(
		ctx context.Context,
		fullShortUrl string,
		statsInfo valobj.ShortLinkStatsRecordVo,
	) (string, error)

	//RecordLinkVisitInfo(ctx context.Context, info valobj.ShortLinkStatsRecordVo) error
}
