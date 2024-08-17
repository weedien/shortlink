package command

import (
	"context"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type SaveToRecycleBinHandler struct {
	repo recycle_bin.Repository
}

func (h SaveToRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
	return h.repo.SaveToRecycleBin(
		ctx,
		id,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			link.Disable()
			return link, nil
		},
	)
}
