package domain

import (
	"context"
	"shortlink/internal/link/domain/link"
)

type LinkRepository interface {
	// ShortUriExists 短链接是否存在
	ShortUriExists(ctx context.Context, shortUrl string) (bool, error)

	// CreateLink 创建短链接
	CreateLink(ctx context.Context, lk *link.Link) error

	// CreateLinkBatch 批量创建短链接
	CreateLinkBatch(ctx context.Context, links []*link.Link) error

	// UpdateLink 更新短链接
	UpdateLink(
		ctx context.Context,
		id link.Identifier,
		updateFn func(ctx context.Context, link *link.Link) (*link.Link, error),
	) error

	// SaveToRecycleBin 保存到回收站
	SaveToRecycleBin(
		ctx context.Context,
		id link.Identifier,
	) error

	// RemoveFromRecycleBin 从回收站移除
	RemoveFromRecycleBin(
		ctx context.Context,
		id link.Identifier,
	) error

	// RecoverFromRecycleBin 恢复回收站
	RecoverFromRecycleBin(
		ctx context.Context,
		id link.Identifier,
	) error
}
