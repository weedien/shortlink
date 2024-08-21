package command

import (
	"context"
	"log/slog"
	"shortlink/common/consts"
	"shortlink/common/decorator"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type recoverFromRecycleBinHandler struct {
	repo recycle_bin.Repository
}

type RecoverFromRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewRecoverFromRecycleBinHandler(
	repo recycle_bin.Repository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RecoverFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[types.LinkID](
		recoverFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h recoverFromRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
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
