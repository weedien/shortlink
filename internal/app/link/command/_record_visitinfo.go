package command

import (
	"context"
	"log/slog"
	"shortlink/common/decorator"
	"shortlink/common/metrics"
	"shortlink/internal/domain/link"
	"shortlink/internal/domain/link/valobj"
)

type recordLinkVisitInfoHandler struct {
	repo link.Repository
}

type RecordLinkVisitInfoHandler decorator.CommandHandler[RecordLinkVisitInfo]

func NewRecordLinkVisitInfoHandler(
	repo link.Repository,
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
