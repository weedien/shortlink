package query

import (
	"context"
	"errors"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/common/types"
)

type pageLinkHandler struct {
	readModel PageLinkReadModel
}

type PageLinkHandler decorator.QueryHandler[PageLink, *types.PageResp[Link]]

func NewPageLinkHandler(
	readModel PageLinkReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) PageLinkHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[PageLink, *types.PageResp[Link]](
		pageLinkHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

type PageLink struct {
	// 分页请求
	types.PageReq
	// 分组ID
	Gid string
	// 排序标识
	// 取值为: todayPv, todayUv, todayUip, totalPv, totalUv, totalUip
	// 默认为 create_time
	OrderTag string
}

func (pl PageLink) validate() error {
	if pl.Gid == "" {
		return errors.New("gid is required")
	}
	return nil
}

type PageLinkReadModel interface {
	PageLink(ctx context.Context, param PageLink) (*types.PageResp[Link], error)
}

func (h pageLinkHandler) Handle(ctx context.Context, param PageLink) (*types.PageResp[Link], error) {
	return h.readModel.PageLink(ctx, param)
}
