package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"shortlink/internal/common/base_event"
	"shortlink/internal/common/lock"
	"shortlink/internal/common/metrics"
	"shortlink/internal/link/adapter"
	"shortlink/internal/link/adapter/readrepo"
	"shortlink/internal/link/app"
	"shortlink/internal/link/app/command"
	"shortlink/internal/link/app/query"
)

func NewShortLinkApplication(
	db *gorm.DB,
	rdb *redis.Client,
	locker lock.DistributedLock,
	eventBus base_event.EventBus,
) (a app.Application) {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	repository := adapter.NewLinkRepository(db, rdb, locker)
	readModel := readrepo.NewLinkQuery(db)

	// 可以理解为一个 CQRS 的总线
	a = app.Application{
		Commands: app.Commands{
			CreateLink:      command.NewCreateLinkHandler(repository, locker, logger, metricsClient),
			CreateLinkBatch: command.NewCreateLinkBatchHandler(repository, logger, metricsClient),
			UpdateLink:      command.NewUpdateLinkHandler(repository, logger, metricsClient),
		},
		Queries: app.Queries{
			PageLink:       query.NewPageLinkHandler(readModel, logger, metricsClient),
			ListGroupCount: query.NewListGroupCountHandler(readModel, logger, metricsClient),
			GetOriginalUrl: query.NewGetOriginalUrlHandler(readModel, eventBus, logger, metricsClient),
		},
	}

	return

}
