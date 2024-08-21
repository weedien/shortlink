package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/common/types"
	"shortlink/internal/recycle_bin/domain"
)

type saveToRecycleBinHandler struct {
	repo domain.Repository
}

type SaveToRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewSaveToRecycleBinHandler(
	repo domain.Repository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) SaveToRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[types.LinkID](
		saveToRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h saveToRecycleBinHandler) Handle(ctx context.Context, id types.LinkID) error {
	return h.repo.SaveToRecycleBin(
		ctx,
		id,
		func(ctx context.Context, link *types.Link) (*types.Link, error) {
			link.Disable()
			return link, nil
		},
	)
}
