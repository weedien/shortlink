package command

import (
	"context"
	"log/slog"
	"shortlink/common/consts"
	"shortlink/common/decorator"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type removeFromRecycleBinHandler struct {
	repo recycle_bin.Repository
}

type RemoveFromRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewRemoveFromRecycleBinHandler(
	repo recycle_bin.Repository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RemoveFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[types.LinkID](
		removeFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h removeFromRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
	return h.repo.RemoveFromRecycleBin(ctx, id, consts.StatusDisable)
}
