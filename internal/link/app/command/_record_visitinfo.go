package command

import (
	"context"
	"log/slog"
	"shortlink/internal/common/decorator"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/domain"
	"shortlink/internal/link/domain/valobj"
)

type recordLinkVisitInfoHandler struct {
	repo domain.LinkRepository
}

type RecordLinkVisitInfoHandler decorator.CommandHandler[RecordLinkVisitInfo]

func NewRecordLinkVisitInfoHandler(
	repo domain.LinkRepository,
	logger *slog.Logger,
	metricsClient metrics.Client,
) RecordLinkVisitInfoHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators(
		recordLinkVisitInfoHandler{repo: repo},
		logger,
		metricsClient,
	)
}

type RecordLinkVisitInfo valobj.ShortLinkStatsRecordVo

func (h recordLinkVisitInfoHandler) Handle(ctx context.Context, cmd RecordLinkVisitInfo) error {
	return h.repo.RecordLinkVisitInfo(ctx, valobj.ShortLinkStatsRecordVo(cmd))
}
