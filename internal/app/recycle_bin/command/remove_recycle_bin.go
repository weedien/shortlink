package command

import (
	"context"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type RemoveFromRecycleBinHandler struct {
	repo recycle_bin.Repository
}

func (h RemoveFromRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
	return h.repo.RemoveFromRecycleBin(ctx, id, consts.StatusDisable)
}
