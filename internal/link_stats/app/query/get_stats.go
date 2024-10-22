package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"time"
)

type getLinkStatHandler struct {
	readModel GetLinkStatReadModel
}

type GetLinkStatHandler decorator.QueryHandler[GetLinkStat, *LinkStat]

func NewGetLinkStatHandler(
	readModel GetLinkStatReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) GetLinkStatHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetLinkStat, *LinkStat](
		getLinkStatHandler{readModel},
		logger,
		metricsClient,
	)
}

type GetLinkStat struct {
	// 完整短链接
	FullShortUrl string
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
	// 启用标识
	EnableStatus int
}

type GetLinkStatReadModel interface {
	// GetLinkStat 获取单个短链接监控数据
	GetLinkStat(ctx context.Context, param GetLinkStat) (res *LinkStat, err error)
}

func (h getLinkStatHandler) Handle(ctx context.Context, q GetLinkStat) (res *LinkStat, err error) {
	return h.readModel.GetLinkStat(ctx, q)
}
