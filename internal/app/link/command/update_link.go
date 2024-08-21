package command

import (
	"context"
	"log/slog"
	"shortlink/common/consts"
	"shortlink/common/decorator"
	"shortlink/common/types"
	"shortlink/internal/domain/link"
	"time"
)

type updateLinkHandler struct {
	repo link.Repository
}

type UpdateLinkHandler decorator.CommandHandler[UpdateLink]

func NewUpdateLinkHandler(
	repo link.Repository,
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
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int
	// 有效期
	ValidDate time.Time
	// 描述
	Description string
}

func (h updateLinkHandler) Handle(ctx context.Context, cmd UpdateLink) (err error) {
	return h.repo.UpdateLink(
		ctx,
		types.LinkID{
			FullShortUrl: cmd.FullShortUrl,
			Gid:          cmd.OriginalGid,
		},
		consts.StatusEnable,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			link.SetGid(cmd.Gid).SetOriginalUrl(cmd.OriginalUrl).
				SetDesc(cmd.Description).SetValidDateType(cmd.ValidDateType).
				SetValidDate(cmd.ValidDate)
			return link, nil
		},
	)
}
