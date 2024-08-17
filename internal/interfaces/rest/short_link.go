package rest

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/internal/app/link"
)

type ShortLinkServer struct {
	app link.Application
}

func NewShortLinkHandler(app link.Application, router fiber.Router) {
	server := &ShortLinkServer{
		app: app,
	}

	// 短链接未找到
	router.All("/page/notfound", func(c *fiber.Ctx) error {
		return c.SendFile(".../resources/notfound.html")
	})
	// 短链接跳转到原始链接
	router.Get("/:short-uri", server.Redirect)
	// 创建短链接
	router.Post("/create", server.CreateShortLink)
	// 通过分布式锁创建短链接
	router.Post("/create/with-lock", server.CreateShortLinkWithLock)
	// 批量创建短链接
	router.Post("/batch-create", server.BatchCreateShortLink)
	// 更新短链接
	router.Put("/update", server.UpdateShortLink)
	// 分页查询短链接
	router.Get("/page", server.PageQueryShortLink)
	// 查询短链接分组内数量
	router.Get("/group-count", server.ListGroupShortLinkCount)
}

// Redirect 短链接跳转到原始链接
func (h ShortLinkServer) Redirect(c *fiber.Ctx) error {
	return nil
}

// CreateShortLink 创建短链接
func (h ShortLinkServer) CreateShortLink(c *fiber.Ctx) error {
	//return h.app.Commands.CreateLink(c)
	return nil
}

// CreateShortLinkWithLock 加锁创建短链接
func (h ShortLinkServer) CreateShortLinkWithLock(c *fiber.Ctx) error {
	return nil
}

// BatchCreateShortLink 批量创建短链接
func (h ShortLinkServer) BatchCreateShortLink(c *fiber.Ctx) error {
	return nil
}

// UpdateShortLink 更新短链接
func (h ShortLinkServer) UpdateShortLink(c *fiber.Ctx) error {
	return nil
}

// PageQueryShortLink 分页查询短链接
func (h ShortLinkServer) PageQueryShortLink(c *fiber.Ctx) error {
	return nil
}

// ListGroupShortLinkCount 查询短链接分组内数量
func (h ShortLinkServer) ListGroupShortLinkCount(c *fiber.Ctx) error {
	return nil
}
