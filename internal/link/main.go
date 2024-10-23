package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"shortlink/internal/common/cache"
	"shortlink/internal/common/config"
	"shortlink/internal/common/database"
	"shortlink/internal/common/lock"
	"shortlink/internal/common/logging"
	"shortlink/internal/common/mq"
	"shortlink/internal/common/server"
	"shortlink/internal/common/shutdown"
	linkservice "shortlink/internal/link/service"
	linktrigger "shortlink/internal/link/trigger/http"
	"syscall"
)

func main() {
	fmt.Println("This is the link-service")

	// 全局日志初始化
	logging.InitLogger()

	// 初始化外部依赖
	db := database.ConnectToDatabase()                            // Postgresql
	rdb := cache.ConnectToRedis()                                 // Redis
	locker := lock.NewRedisLock(rdb)                              // DistributedLock - Redis
	eventBus := mq.NewRocketMqBasedEventBus(context.Background()) // EventBus

	// 创建应用服务
	shortLinkApp := linkservice.NewLinkApplication(db, rdb, locker, eventBus)

	shutdownServer := server.RunHttpServerOnPort(config.Port.String(), func(router fiber.Router) {
		server.NewUriTitleApi(router)
		linktrigger.NewLinkApi(shortLinkApp, router)
	})

	shutdown.NewHook().WithSignals(syscall.SIGINT, syscall.SIGTERM).Close(
		// shutdown server
		shutdownServer,
		// shutdown database
		func() {
			if sqlDB, err := db.DB(); err != nil {
				slog.Error("database.DB() failed", "error", err)
			} else {
				if err = sqlDB.Close(); err != nil {
					slog.Error("database.DB().Close() failed", "error", err)
				}
			}
		},
		// shutdown redis
		func() {
			if err := rdb.Close(); err != nil {
				slog.Error("redis.Close() failed", "error", err)
			}
		},
		// shutdown event bus
		func() {
			if eventBus != nil {
				// 事件总线是自己封装的，关闭失败的情况已经在内部进行了处理
				slog.Info("Closing event bus")
				eventBus.Close()
			}
		},
	)
}
