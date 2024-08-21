package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"shortlink/common/metrics"
	"shortlink/internal/app/recycle_bin"
	"shortlink/internal/app/recycle_bin/command"
	"shortlink/internal/app/recycle_bin/query"
	"shortlink/internal/infra/persistence/readrepo"
	"shortlink/internal/infra/persistence/repo"
)

func NewShortLinkRecycleBinApplication(
	db *gorm.DB,
	rdb *redis.Client,
) recycle_bin.Application {

	logger := slog.Default()
	metricsClient := metrics.NoOp{}
	repository := repo.NewRecycleBinRepository(db, rdb)
	readModel := readrepo.NewRecycleBinQuery(db)

	return recycle_bin.Application{
		Commands: recycle_bin.Commands{
			SaveToRecycleBin:      command.NewSaveToRecycleBinHandler(repository, logger, metricsClient),
			RemoveFromRecycleBin:  command.NewRemoveFromRecycleBinHandler(repository, logger, metricsClient),
			RecoverFromRecycleBin: command.NewRecoverFromRecycleBinHandler(repository, logger, metricsClient),
		},
		Queries: recycle_bin.Queries{
			PageDisabledLink: query.NewPageRecycleBinHandler(readModel, logger, metricsClient),
		},
	}
}
