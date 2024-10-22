package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/link"
)

type recoverFromRecycleBinHandler struct {
	repo domain.LinkRepository
}

type RecoverFromRecycleBinHandler decorator.CommandHandler[link.Identifier]

func NewRecoverFromRecycleBinHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RecoverFromRecycleBinHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[link.Identifier](
		recoverFromRecycleBinHandler{repo},
		logger,
		metricsClient,
	)
}

func (h recoverFromRecycleBinHandler) Handle(ctx context.Context, id link.Identifier) error {
	return h.repo.RecoverFromRecycleBin(ctx, id)
}
