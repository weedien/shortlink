package service

import (
	"gorm.io/gorm"
	"log/slog"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link_stats/app"
	"shortlink/internal/link_stats/app/query"
	"shortlink/internal/link_stats/repo/readrepo"
)

func NewShortLinkStatsApplication(db *gorm.DB) app.Application {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	readModel := readrepo.NewLinkStatsQuery(db)

	return app.Application{
		Queries: app.Queries{
			GetLinkStats:               query.NewGetLinkStatsHandler(readModel, logger, metricsClient),
			GroupLinkStats:             query.NewGroupLinkStatsHandler(readModel, logger, metricsClient),
			GetLinkStatsAccessRecord:   query.NewGetLinkStatsAccessRecordHandler(readModel, logger, metricsClient),
			GroupLinkStatsAccessRecord: query.NewGroupLinkStatsAccessRecordHandler(readModel, logger, metricsClient),
		},
	}
}
