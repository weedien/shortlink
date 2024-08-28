package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"log/slog"
	"shortlink/internal/common/config"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/server/httperr"
)

func RunHttpServerOnPort(port string, createHandler func(router fiber.Router)) func() {
	app := fiber.New(fiber.Config{
		//AppName:      "short-link-service",
		ErrorHandler: httperr.ErrorHandler,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
	})
	setupMiddlewares(app)

	// 监控页面
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics Page"}))

	apiRouter := app.Group(config.BaseRoutePrefix.String())
	createHandler(apiRouter)

	// 处理未找到的路由
	app.All("*", func(c *fiber.Ctx) error {
		return error_no.RouteNotFound
	})

	slog.Info("HTTP server is running on port " + port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			slog.Error("HTTP server failed to start", "error", err)
		}
	}()

	return func() {
		slog.Info("HTTP server is shutting down")
		if err := app.Shutdown(); err != nil {
			slog.Error("HTTP server failed to shut down", "error", err)
		}
	}
}
