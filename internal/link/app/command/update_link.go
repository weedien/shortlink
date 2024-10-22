package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/link"
	"time"
)

type updateLinkHandler struct {
	repo domain.LinkRepository
}

type UpdateLinkHandler decorator.CommandHandler[UpdateLink]

func NewUpdateLinkHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) UpdateLinkHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[UpdateLink](
		updateLinkHandler{repo: repo},
		logger,
		metricsClient,
	)
}

type UpdateLink struct {
	// 完整短链接
	FullShortUrl string
	// 原始链接
	OriginalUrl string
	// 原始分组ID
	OriginalGid string
	// 分组ID
	Gid string
	// 状态
	Status string
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidType int
	// 有效期 - 开始时间
	ValidStartDate time.Time
	// 有效期 - 结束时间
	ValidEndDate time.Time
	// 描述
	Desc string
}

func (h updateLinkHandler) Handle(ctx context.Context, cmd UpdateLink) (err error) {
	return h.repo.UpdateLink(
		ctx,
		link.Identifier{
			ShortUri: cmd.FullShortUrl,
			Gid:      cmd.OriginalGid,
		},
		func(ctx context.Context, lk *link.Link) (*link.Link, error) {
			err = lk.Update(&cmd.Gid, &cmd.OriginalUrl, &cmd.Status, &cmd.ValidType, &cmd.ValidEndDate, &cmd.Desc)
			if err != nil {
				return nil, err
			}
			return lk, nil
		},
	)
}
