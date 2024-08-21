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

type removeFromRecycleBinHandler struct {
	repo domain.Repository
}

type RemoveFromRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewRemoveFromRecycleBinHandler(
	repo domain.Repository,
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
	return h.repo.RemoveFromRecycleBin(ctx, id, constant.StatusDisable)
}
