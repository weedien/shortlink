package service

import (
	"gorm.io/gorm"
	"log/slog"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link_stats/app"
	"shortlink/internal/link_stats/app/query"
	"shortlink/internal/link_stats/repo/readrepo"
)

func NewShortLinkStatApplication(db *gorm.DB) app.Application {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	readModel := readrepo.NewLinkStatQuery(db)

	return app.Application{
		Queries: app.Queries{
			GetLinkStat:               query.NewGetLinkStatHandler(readModel, logger, metricsClient),
			GroupLinkStat:             query.NewGroupLinkStatHandler(readModel, logger, metricsClient),
			GetLinkStatAccessRecord:   query.NewGetLinkStatAccessRecordHandler(readModel, logger, metricsClient),
			GroupLinkStatAccessRecord: query.NewGroupLinkStatAccessRecordHandler(readModel, logger, metricsClient),
		},
	}
}
