package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/app/recycle_bin"
	"shortlink/internal/app/recycle_bin/query"
	"shortlink/internal/interfaces/rest/dto/req"
	"shortlink/internal/interfaces/rest/dto/resp"
)

type RecycleBinApi struct {
	app recycle_bin.Application
}

func NewShortLinkRecycleBinApi(app recycle_bin.Application, router fiber.Router) {
	api := &RecycleBinApi{
		app: app,
	}

	recycleBin := router.Group("/recycle-bin")
	// 保存到回收站
	recycleBin.Post("/recycle-bin/save", api.SaveToRecycleBin)
	// 分页查询回收站短链接
	recycleBin.Get("/recycle-bin/page", api.PageQueryRecycleBin)
	// 恢复短链接
	recycleBin.Post("/recycle-bin/recover", api.RecoverShortLink)
	// 从回收站移除短链接
	recycleBin.Delete("/recycle-bin/remove", api.RemoveFromRecycleBin)
}

// SaveToRecycleBin 保存到回收站
func (h RecycleBinApi) SaveToRecycleBin(c *fiber.Ctx) error {
	var reqParam req.RecycleBinSaveReq
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	err := h.app.Commands.SaveToRecycleBin.Handle(c.Context(), types.LinkID{
		Gid:          reqParam.Gid,
		FullShortUrl: reqParam.FullShortUrl,
	})
	if err != nil {
		return err
	}

	c.Status(fiber.StatusNoContent)

	return nil
}

// PageQueryRecycleBin 分页查询回收站短链接
func (h RecycleBinApi) PageQueryRecycleBin(c *fiber.Ctx) error {
	var reqParam req.RecycleBinPageReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.PageDisabledLink.Handle(c.Context(), query.PageRecycleBin{
		PageReq:      reqParam.PageReq,
		GidList:      reqParam.GidList,
		EnableStatus: consts.StatusDisable,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkPageResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// RecoverShortLink 恢复短链接
func (h RecycleBinApi) RecoverShortLink(c *fiber.Ctx) error {
	var reqParam req.RecycleBinRecoverReq
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	err := h.app.Commands.RecoverFromRecycleBin.Handle(c.Context(), types.LinkID{
		Gid:          reqParam.Gid,
		FullShortUrl: reqParam.FullShortUrl,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	c.Status(fiber.StatusNoContent)

	return nil
}

// RemoveFromRecycleBin 从回收站移除短链接
func (h RecycleBinApi) RemoveFromRecycleBin(c *fiber.Ctx) error {
	var reqParam req.RecycleBinDeleteReq
	if err := c.BodyParser(&reqParam); err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	err := h.app.Commands.RemoveFromRecycleBin.Handle(c.Context(), types.LinkID{
		Gid:          reqParam.Gid,
		FullShortUrl: reqParam.FullShortUrl,
	})
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	c.Status(fiber.StatusNoContent)

	return nil
}
