package service

import "shortlink/internal/infra/persistence/repo"

// LinkAppService 短链接应用服务，
// 将多个handler的共用逻辑抽取到应用服务中
type LinkAppService struct {
	repo repo.LinkRepository
}

func (s LinkAppService) CreateLink() error {
	return nil
}

func (s LinkAppService) CreateLinkWithLock() error {
	return nil
}
