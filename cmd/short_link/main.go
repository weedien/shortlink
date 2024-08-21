package main

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"shortlink/common/config"
	"shortlink/common/error_no"
	"shortlink/common/httperr"
	"shortlink/common/logging"
	"shortlink/internal/infra/cache"
	"shortlink/internal/infra/lock"
	"shortlink/internal/infra/mq"
	"shortlink/internal/infra/persistence/database"
	"shortlink/internal/interfaces/middleware"
	"shortlink/internal/interfaces/rest"
	"shortlink/internal/service"
	"syscall"
)

func main() {
	// 全局日志初始化
	logging.InitLogger()

	// 初始化外部依赖
	db := database.ConnectToDatabase()                            // Postgresql
	rdb := cache.ConnectToRedis()                                 // Redis
	locker := lock.NewRedisLock(rdb)                              // DistributedLock - Redis
	eventBus := mq.NewRocketMqBasedEventBus(context.Background()) // EventBus

	// 创建 Fiber 应用
	f := fiber.New(fiber.Config{
		AppName:      "short-link-service",
		ErrorHandler: httperr.ErrorHandler,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})

	// 注册全局中间件
	middleware.SetupMiddlewares(f)

	// 创建应用服务
	shortLinkApp := service.NewShortLinkApplication(db, rdb, locker, eventBus)
	RecycleBinApp := service.NewShortLinkRecycleBinApplication(db, rdb)
	shortLinkStatsApp := service.NewShortLinkStatsApplication(db)

	// 注册路由
	router := f.Group(config.BaseRoutePrefix.String())
	rest.NewUriTitleApi(router)
	rest.NewShortLinkApi(shortLinkApp, router)
	rest.NewShortLinkStatsApi(shortLinkStatsApp, router)
	rest.NewShortLinkRecycleBinApi(RecycleBinApp, router)

	// 处理未找到的路由
	f.All("*", func(c *fiber.Ctx) error {
		return error_no.RouteNotFound
	})

	f.Get("/metrics", monitor.New(monitor.Config{Title: "ShortLinkService Metrics Page"}))

	go func() {
		log.Fatal(f.Listen(":" + config.Port.String()))
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = f.Shutdown()

	fmt.Println("Running cleanup tasks...")
	// shutdown database
	if sqlDB, err := db.DB(); err != nil {
		slog.Warn("database.DB() failed", "error", err)
	} else {
		err = sqlDB.Close()
		if err != nil {
			slog.Warn("database.DB().Close() failed", "error", err)
		}
	}
	// shutdown redis
	if err := rdb.Close(); err != nil {
		slog.Warn("redis gracefully shutdown failed", "error", err)
	}
	// shutdown event bus
	eventBus.Close() // 事件总线是自己封装的，关闭失败的情况已经在内部进行了处理
	slog.Info("short-link-service was successful shutdown.")
}
