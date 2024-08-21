package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"shortlink/internal/common/metrics"
	"shortlink/internal/recycle_bin/adapters"
	"shortlink/internal/recycle_bin/adapters/readrepo"
	"shortlink/internal/recycle_bin/app"
	"shortlink/internal/recycle_bin/app/command"
	"shortlink/internal/recycle_bin/app/query"
)

func NewShortLinkRecycleBinApplication(
	db *gorm.DB,
	rdb *redis.Client,
) app.Application {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	repository := adapters.NewRecycleBinRepository(db, rdb)
	readModel := readrepo.NewRecycleBinQuery(db)

	return app.Application{
		Commands: app.Commands{
			SaveToRecycleBin:      command.NewSaveToRecycleBinHandler(repository, logger, metricsClient),
			RemoveFromRecycleBin:  command.NewRemoveFromRecycleBinHandler(repository, logger, metricsClient),
			RecoverFromRecycleBin: command.NewRecoverFromRecycleBinHandler(repository, logger, metricsClient),
		},
		Queries: app.Queries{
			PageDisabledLink: query.NewPageRecycleBinHandler(readModel, logger, metricsClient),
		},
	}
}
