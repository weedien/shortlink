package recycle_bin

import (
	"context"
	"shortlink/common/types"
)

type Repository interface {

	// SaveToRecycleBin 保存到回收站
	SaveToRecycleBin(
		ctx context.Context,
		id types.LinkID,
		updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
	) error

	// RemoveFromRecycleBin 从回收站移除
	RemoveFromRecycleBin(
		ctx context.Context,
		id types.LinkID,
		enableStatus int,
	) error

	// RecoverFromRecycleBin 恢复回收站
	RecoverFromRecycleBin(
		ctx context.Context,
		id types.LinkID,
		enableStatus int,
		updateFn func(ctx context.Context, link *types.Link) (*types.Link, error),
	) error
}
