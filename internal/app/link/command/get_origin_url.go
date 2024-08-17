package command

import (
	"context"
	"shortlink/internal/domain/link"
	"shortlink/internal/domain/link/valobj"
)

type GetOriginUrlHandler struct {
	repo link.Repository
}

type GetOriginUrl struct {
	FullShortUrl string
	StatsInfo    valobj.ShortLinkStatsRecordVO
}

func (h GetOriginUrlHandler) Handle(ctx context.Context, cmd GetOriginUrl) (string, error) {
	return h.repo.GetOriginalUrlByShortUrl(ctx, cmd.FullShortUrl, cmd.StatsInfo)
}
