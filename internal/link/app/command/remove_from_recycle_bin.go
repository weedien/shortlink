package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/link"
)

type removeFromRecycleBinHandler struct {
	repo domain.LinkRepository
}

type RemoveFromRecycleBinHandler decorator.CommandHandler[link.Identifier]

func NewRemoveFromRecycleBinHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RemoveFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[link.Identifier](
		removeFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h removeFromRecycleBinHandler) Handle(ctx context.Context, id link.Identifier) error {
	return h.repo.RemoveFromRecycleBin(ctx, id)
}
