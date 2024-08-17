package query

import (
	"context"
	"shortlink/common/types"
	"time"
)

// GetLinkStatsAccessRecordHandler 获取单个短链接指定时间内访问记录监控数据
type GetLinkStatsAccessRecordHandler struct {
	readModel GetLinkStatsAccessRecordReadModel
}

type GetLinkStatsAccessRecordGroup struct {
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

type GetLinkStatsAccessRecordReadModel interface {
	// GetLinkStatsAccessRecord 获取单个短链接指定时间内访问记录监控数据
	GetLinkStatsAccessRecord(ctx context.Context, param GetLinkStatsAccessRecordGroup) (*types.PageResp[LinkStatsAccessRecord], error)
}

func (h GetLinkStatsAccessRecordHandler) Handle(ctx context.Context, query GetLinkStatsAccessRecordGroup) (d *types.PageResp[LinkStatsAccessRecord], err error) {
	return h.readModel.GetLinkStatsAccessRecord(ctx, query)
}
