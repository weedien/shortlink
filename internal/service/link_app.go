package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"reflect"
	"shortlink/common/base_event"
	"shortlink/common/metrics"
	"shortlink/internal/app/link"
	"shortlink/internal/app/link/command"
	"shortlink/internal/app/link/event"
	"shortlink/internal/app/link/query"
	appservice "shortlink/internal/app/link/service"
	"shortlink/internal/app/link_stats/listener"
	"shortlink/internal/infra/lock"
	"shortlink/internal/infra/persistence/readrepo"
	"shortlink/internal/infra/persistence/repo"
)

func NewShortLinkApplication(
	db *gorm.DB,
	rdb *redis.Client,
	locker lock.DistributedLock,
	eventBus base_event.AppEventBus,
) (app link.Application) {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	repository := repo.NewLinkRepository(db, rdb, locker)
	readModel := readrepo.NewLinkQuery(db)

	// 可以理解为一个 CQRS 的总线
	app = link.Application{
		Commands: link.Commands{
			CreateLink:      command.NewCreateLinkHandler(repository, locker, logger, metricsClient),
			CreateLinkBatch: command.NewCreateLinkBatchHandler(repository, logger, metricsClient),
			UpdateLink:      command.NewUpdateLinkHandler(repository, logger, metricsClient),
		},
		Queries: link.Queries{
			PageLink:       query.NewPageLinkHandler(readModel, logger, metricsClient),
			ListGroupCount: query.NewListGroupCountHandler(readModel, logger, metricsClient),
			GetOriginalUrl: query.NewGetOriginalUrlHandler(readModel, eventBus, logger, metricsClient),
		},
	}

	// 应用服务
	service := appservice.NewLinkAppService(repository)

	// $$ 消费者订阅事件
	linkAccessedListener := listener.NewLinkAccessedEventListener(service)
	eventBus.Subscribe(reflect.TypeOf(event.RecordLinkVisitEvent{}), linkAccessedListener)

	return

}
