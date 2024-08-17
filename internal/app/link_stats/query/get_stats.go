package query

import (
	"context"
	"time"
)

type GetLinkStatsHandler struct {
	readModel GetLinkStatsReadModel
}

type GetLinkStats struct {
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

type GetLinkStatsReadModel interface {
	// GetLinkStats 获取单个短链接监控数据
	GetLinkStats(ctx context.Context, param GetLinkStats) (res *LinkStats, err error)
}

func (h GetLinkStatsHandler) Handle(ctx context.Context, q GetLinkStats) (res *LinkStats, err error) {
	return h.readModel.GetLinkStats(ctx, q)
}
