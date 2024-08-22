package domain

import (
	"context"
	"shortlink/internal/link/domain/entity"
)

type RecycleBinRepository interface {

	// SaveToRecycleBin 保存到回收站
	SaveToRecycleBin(
		ctx context.Context,
		id entity.LinkID,
		updateFn func(ctx context.Context, link *entity.Link) (*entity.Link, error),
	) error

	// RemoveFromRecycleBin 从回收站移除
	RemoveFromRecycleBin(
		ctx context.Context,
		id entity.LinkID,
		enableStatus int,
	) error

	// RecoverFromRecycleBin 恢复回收站
	RecoverFromRecycleBin(
		ctx context.Context,
		id entity.LinkID,
		enableStatus int,
		updateFn func(ctx context.Context, link *entity.Link) (*entity.Link, error),
	) error
}
