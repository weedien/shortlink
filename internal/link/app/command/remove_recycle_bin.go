package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/entity"
)

type removeFromRecycleBinHandler struct {
	repo domain.RecycleBinRepository
}

type RemoveFromRecycleBinHandler decorator.CommandHandler[entity.LinkID]

func NewRemoveFromRecycleBinHandler(
	repo domain.RecycleBinRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RemoveFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[entity.LinkID](
		removeFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h removeFromRecycleBinHandler) Handle(ctx context.Context, id entity.LinkID) error {
	return h.repo.RemoveFromRecycleBin(ctx, id, constant.StatusDisable)
}
