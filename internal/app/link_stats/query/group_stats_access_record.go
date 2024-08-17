package query

import (
	"context"
	"shortlink/common/types"
	"time"
)

type GroupLinkStatsAccessRecordHandler struct {
	readModel GroupLinkStatsAccessRecordReadModel
}

type GroupLinkStatsAccessRecord struct {
	// 分页参数
	types.PageReq
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
}

type GroupLinkStatsAccessRecordReadModel interface {
	// GroupLinkStatsAccessRecord 获取分组指定时间内访问记录监控数据
	GroupLinkStatsAccessRecord(ctx context.Context, param GroupLinkStatsAccessRecord) (*types.PageResp[LinkStatsAccessRecord], error)
}

func (h GroupLinkStatsAccessRecordHandler) Handle(ctx context.Context, q GroupLinkStatsAccessRecord) (*types.PageResp[LinkStatsAccessRecord], error) {
	return h.readModel.GroupLinkStatsAccessRecord(ctx, q)
}
