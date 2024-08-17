package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"shortlink/config"
	"shortlink/internal/app/link"
	"shortlink/internal/app/recycle_bin"
	"shortlink/internal/infra/cache"
	"shortlink/internal/infra/persistence/database"
	"shortlink/internal/ports/middleware"
	"shortlink/pkg/error_no"
	"shortlink/pkg/types"
)

func main() {
	// 连接数据库
	database.ConnectToDatabase()
	cache.ConnectToRedis()

	// 创建 Fiber 应用
	f := fiber.New(fiber.Config{
		AppName:      "Link",
		ServerHeader: "Link",
	})

	// 注册全局中间件
	middleware.SetupMiddlewares(f)

	// 创建应用服务
	shortLinkApp := link.NewLinkApplication()
	RecycleBinApp := recycle_bin.NewRecycleBinApplication()

	// 注册路由
	router := f.Group(config.Default("BASE_ROUTE_PREFIX", config.BaseRoutePrefix))
	rest.NewUriTitleServer(router)
	rest.NewShortLinkHandler(shortLinkApp, router)
	rest.NewRecycleBinHandler(RecycleBinApp, router)

	// 处理未找到的路由
	f.All("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(
			types.FailWithErrorCode(error_no.RouteNotFound),
		)
	})

	log.Fatal(f.Listen(config.Default("PORT", config.PORT)))
}
