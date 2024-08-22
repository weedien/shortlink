package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/entity"
)

type saveToRecycleBinHandler struct {
	repo domain.RecycleBinRepository
}

type SaveToRecycleBinHandler decorator.CommandHandler[entity.LinkID]

func NewSaveToRecycleBinHandler(
	repo domain.RecycleBinRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) SaveToRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[entity.LinkID](
		saveToRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h saveToRecycleBinHandler) Handle(ctx context.Context, id entity.LinkID) error {
	return h.repo.SaveToRecycleBin(
		ctx,
		id,
		func(ctx context.Context, link *entity.Link) (*entity.Link, error) {
			link.Disable()
			return link, nil
		},
	)
}
