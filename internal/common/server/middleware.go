package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"shortlink/internal/common/config"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/server/httperr"
)

func setupMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(limiter.New(limiter.Config{
		Max: config.DefaultInt("MAX_REQUESTS", 100),
		LimitReached: func(c *fiber.Ctx) error {
			return httperr.RespondWithError(c, error_no.TooManyRequests)
		},
	}))
	app.Use(logger.New())
	//app.Use(recover.New())
	app.Use(requestid.New())
}
