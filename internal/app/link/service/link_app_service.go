package service

import (
	"context"
	"shortlink/internal/domain/link/valobj"
	"shortlink/internal/infra/persistence/repo"
)

// LinkAppService 短链接应用服务接口
type LinkAppService interface {
	RecordLinkVisitInfo(ctx context.Context, info valobj.ShortLinkStatsRecordVo) error
}

// LinkAppServiceImpl 短链接应用服务实现
type LinkAppServiceImpl struct {
	repo repo.LinkRepository
}

func NewLinkAppService(repo repo.LinkRepository) LinkAppService {
	return LinkAppServiceImpl{repo: repo}
}

// RecordLinkVisitInfo 记录短链接访问信息 这个接口由应用事件触发，作为用户请求短链接的副作用
func (s LinkAppServiceImpl) RecordLinkVisitInfo(ctx context.Context, info valobj.ShortLinkStatsRecordVo) error {
	return s.repo.RecordLinkVisitInfo(ctx, info)
}
