package query

import (
	"context"
	"time"
)

type GroupLinkStatsHandler struct {
	readModel GroupLinkStatsReadModel
}

type GroupLinkStats struct {
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
}

type GroupLinkStatsReadModel interface {
	// GroupLinkStats 获取分组短链接监控数据
	GroupLinkStats(ctx context.Context, param GroupLinkStats) (*LinkStats, error)
}

func (h GroupLinkStatsHandler) Handle(ctx context.Context, q GroupLinkStats) (*LinkStats, error) {
	return h.readModel.GroupLinkStats(ctx, q)
}
