package service

import (
	"gorm.io/gorm"
	"log/slog"
	"shortlink/common/metrics"
	"shortlink/internal/app/link_stats"
	"shortlink/internal/app/link_stats/query"
	"shortlink/internal/infra/persistence/readrepo"
)

func NewShortLinkStatsApplication(db *gorm.DB) link_stats.Application {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	readModel := readrepo.NewLinkStatsQuery(db)

	return link_stats.Application{
		Queries: link_stats.Queries{
			GetLinkStats:               query.NewGetLinkStatsHandler(readModel, logger, metricsClient),
			GroupLinkStats:             query.NewGroupLinkStatsHandler(readModel, logger, metricsClient),
			GetLinkStatsAccessRecord:   query.NewGetLinkStatsAccessRecordHandler(readModel, logger, metricsClient),
			GroupLinkStatsAccessRecord: query.NewGroupLinkStatsAccessRecordHandler(readModel, logger, metricsClient),
		},
	}
}
