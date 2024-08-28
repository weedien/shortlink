package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"os/signal"
	"shortlink/internal/common/cache"
	"shortlink/internal/common/database"
	"shortlink/internal/common/logging"
	"shortlink/internal/common/server"
	"shortlink/internal/common/server/middleware/auth"
	"shortlink/internal/user/service"
	"shortlink/internal/user/trigger/rest"
	"syscall"
)

func main() {
	// 全局日志初始化
	logging.InitLogger()

	// 初始化外部依赖
	db := database.ConnectToDatabase()
	rdb := cache.ConnectToRedis()

	// 创建应用服务
	groupApp := service.NewGroupApplication(db, nil)
	userApp := service.NewUserApplication(db, rdb, groupApp.Commands.CreateGroup)

	cleanup := server.RunHttpServerOnPort("8080", func(router fiber.Router) {
		router.Use(auth.New(rdb)) // 鉴权中间件
		server.NewUriTitleApi(router)
		rest.NewUserApi(userApp, router)
		rest.NewGroupApi(groupApp, router)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	cleanup()

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
	slog.Info("short-link-service was successful shutdown.")

}
