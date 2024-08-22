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

type recoverFromRecycleBinHandler struct {
	repo domain.RecycleBinRepository
}

type RecoverFromRecycleBinHandler decorator.CommandHandler[entity.LinkID]

func NewRecoverFromRecycleBinHandler(
	repo domain.RecycleBinRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RecoverFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[entity.LinkID](
		recoverFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h recoverFromRecycleBinHandler) Handle(ctx context.Context, id entity.LinkID) error {
	return h.repo.RecoverFromRecycleBin(
		ctx,
		id,
		constant.StatusDisable,
		func(ctx context.Context, link *entity.Link) (*entity.Link, error) {
			link.Enable()
			return link, nil
		},
	)
}
