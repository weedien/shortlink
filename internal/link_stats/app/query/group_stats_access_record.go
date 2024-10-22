package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/common/types"
	"time"
)

type groupLinkStatAccessRecordHandler struct {
	readModel GroupLinkStatAccessRecordReadModel
}

type GroupLinkStatAccessRecordHandler decorator.QueryHandler[GroupLinkStatAccessRecord, *types.PageResp[LinkStatAccessRecord]]

func NewGroupLinkStatAccessRecordHandler(
	readModel GroupLinkStatAccessRecordReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) GroupLinkStatAccessRecordHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GroupLinkStatAccessRecord, *types.PageResp[LinkStatAccessRecord]](
		groupLinkStatAccessRecordHandler{readModel},
		logger,
		metricsClient,
	)
}

type GroupLinkStatAccessRecord struct {
	// 分页参数
	types.PageReq
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
}

type GroupLinkStatAccessRecordReadModel interface {
	// GroupLinkStatAccessRecord 获取分组指定时间内访问记录监控数据
	GroupLinkStatAccessRecord(ctx context.Context, param GroupLinkStatAccessRecord) (*types.PageResp[LinkStatAccessRecord], error)
}

func (h groupLinkStatAccessRecordHandler) Handle(ctx context.Context, q GroupLinkStatAccessRecord) (*types.PageResp[LinkStatAccessRecord], error) {
	return h.readModel.GroupLinkStatAccessRecord(ctx, q)
}
