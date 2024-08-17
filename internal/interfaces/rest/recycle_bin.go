package rest

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/internal/app/recycle_bin"
)

type RecycleBinServer struct {
	app recycle_bin.RecycleBinApplication
}

func NewRecycleBinHandler(app recycle_bin.RecycleBinApplication, router fiber.Router) {
	server := &RecycleBinServer{
		app: app,
	}

	recycleBin := router.Group("/recycle-bin")
	// 保存到回收站
	recycleBin.Post("/recycle-bin/save", server.SaveToRecycleBin)
	// 分页查询回收站短链接
	recycleBin.Get("/recycle-bin/page", server.PageQueryRecycleBin)
	// 恢复短链接
	recycleBin.Post("/recycle-bin/recover", server.RecoverShortLink)
	// 从回收站移除短链接
	recycleBin.Delete("/recycle-bin/remove", server.RemoveFromRecycleBin)
}

// SaveToRecycleBin 保存到回收站
func (h RecycleBinServer) SaveToRecycleBin(c *fiber.Ctx) error {
	return nil
}

// PageQueryRecycleBin 分页查询回收站短链接
func (h RecycleBinServer) PageQueryRecycleBin(c *fiber.Ctx) error {
	return nil
}

// RecoverShortLink 恢复短链接
func (h RecycleBinServer) RecoverShortLink(c *fiber.Ctx) error {
	return nil
}

// RemoveFromRecycleBin 从回收站移除短链接
func (h RecycleBinServer) RemoveFromRecycleBin(c *fiber.Ctx) error {
	return nil
}
