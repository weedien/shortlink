package command

import (
	"context"
	"log/slog"
	"shortlink/common/decorator"
	"shortlink/common/types"
	"shortlink/internal/domain/recycle_bin"
)

type saveToRecycleBinHandler struct {
	repo recycle_bin.Repository
}

type SaveToRecycleBinHandler decorator.CommandHandler[types.LinkID]

func NewSaveToRecycleBinHandler(
	repo recycle_bin.Repository,
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
