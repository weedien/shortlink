package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/common/types"
	"shortlink/internal/recycle_bin/domain"
)

type recoverFromRecycleBinHandler struct {
	repo domain.Repository
}

type RecoverFromRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewRecoverFromRecycleBinHandler(
	repo domain.Repository,
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
		constant.StatusDisable,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			link.Enable()
			return link, nil
		},
	)
}
