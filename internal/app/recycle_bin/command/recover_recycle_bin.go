package command

import (
	"context"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type RecoverFromRecycleBinHandler struct {
	repo recycle_bin.Repository
}

func (h RecoverFromRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
	return h.repo.RecoverFromRecycleBin(
		ctx,
		id,
		consts.StatusDisable,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			link.Enable()
			return link, nil
		},
	)
}
