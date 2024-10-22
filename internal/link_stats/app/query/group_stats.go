package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"time"
)

type groupLinkStatHandler struct {
	readModel GroupLinkStatReadModel
}

type GroupLinkStatHandler decorator.QueryHandler[GroupLinkStat, *LinkStat]

func NewGroupLinkStatHandler(
	readModel GroupLinkStatReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) GroupLinkStatHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GroupLinkStat, *LinkStat](
		groupLinkStatHandler{readModel},
		logger,
		metricsClient,
	)
}

type GroupLinkStat struct {
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
}

type GroupLinkStatReadModel interface {
	// GroupLinkStat 获取分组短链接监控数据
	GroupLinkStat(ctx context.Context, param GroupLinkStat) (*LinkStat, error)
}

func (h groupLinkStatHandler) Handle(ctx context.Context, q GroupLinkStat) (*LinkStat, error) {
	return h.readModel.GroupLinkStat(ctx, q)
}
