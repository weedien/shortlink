package query

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/common/types"
	"time"
)

// getLinkStatAccessRecordHandler 获取单个短链接指定时间内访问记录监控数据
type getLinkStatAccessRecordHandler struct {
	readModel GetLinkStatAccessRecordReadModel
}

type GetLinkStatAccessRecordHandler decorator.QueryHandler[GetLinkStatAccessRecord, *types.PageResp[LinkStatAccessRecord]]

func NewGetLinkStatAccessRecordHandler(
	readModel GetLinkStatAccessRecordReadModel,
	logger *slog.Logger,
	metricsClient metrics.Client,
) GetLinkStatAccessRecordHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators[GetLinkStatAccessRecord, *types.PageResp[LinkStatAccessRecord]](
		getLinkStatAccessRecordHandler{readModel},
		logger,
		metricsClient,
	)
}

type GetLinkStatAccessRecord struct {
	types.PageReq
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

type GetLinkStatAccessRecordReadModel interface {
	// GetLinkStatAccessRecord 获取单个短链接指定时间内访问记录监控数据
	GetLinkStatAccessRecord(ctx context.Context, param GetLinkStatAccessRecord) (*types.PageResp[LinkStatAccessRecord], error)
}

func (h getLinkStatAccessRecordHandler) Handle(ctx context.Context, query GetLinkStatAccessRecord) (d *types.PageResp[LinkStatAccessRecord], err error) {
	return h.readModel.GetLinkStatAccessRecord(ctx, query)
}
