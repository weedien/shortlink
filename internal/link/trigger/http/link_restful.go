package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/toolkit"
	"shortlink/internal/link/app"
	"shortlink/internal/link/app/command"
	"shortlink/internal/link/app/query"
	"shortlink/internal/link/domain/valobj"
	"shortlink/internal/link/trigger/http/dto/req"
	"shortlink/internal/link/trigger/http/dto/resp"
	"strings"
	"time"
)

type ShortLinkApi struct {
	app app.Application
}

func NewShortLinkApi(app app.Application, router fiber.Router) {
	api := &ShortLinkApi{
		app: app,
	}

	// 短链接未找到
	router.All("/page/notfound", func(c *fiber.Ctx) error {
		return c.SendFile("resources/notfound.html")
	})
	// 短链接跳转到原始链接
	router.Get("/:short-uri", api.Redirect)
	// 创建短链接
	router.Post("/create", api.CreateShortLink)
	// 通过分布式锁创建短链接
	router.Post("/create/with-lock", api.CreateShortLinkWithLock)
	// 批量创建短链接
	router.Post("/batch-create", api.BatchCreateShortLink)
	// 更新短链接
	router.Put("/update", api.UpdateShortLink)
	// 分页查询短链接
	router.Get("/page", api.PageQueryShortLink)
	// 查询短链接分组内数量
	router.Get("/group-count", api.ListGroupShortLinkCount)
}

// Redirect 短链接跳转到原始链接
func (h ShortLinkApi) Redirect(c *fiber.Ctx) error {

	shortUri := c.Params("short-uri")
	if shortUri == "" {
		return error_no.ShortLinkNotFound
	}

	os, browser, device, network := toolkit.GetRequestInfo(c)

	uv := c.Cookies("uv")
	if uv == "" {
		uv = uuid.NewString()
		c.Cookie(&fiber.Cookie{
			Name:    "uv",
			Value:   uv,
			Expires: time.Now().AddDate(0, 1, 0),
		})
	}

	recordInfo := valobj.ShortLinkStatsRecordVo{
		RemoteAddr:  c.IP(),
		OS:          os,
		Browser:     browser,
		Device:      device,
		Network:     network,
		UV:          uv,
		CurrentDate: time.Now(),
	}

	q := query.GetOriginalUrl{
		FullShortUrl: fmt.Sprintf("https://%s/%s", c.BaseUrl(), shortUri),
		RecordInfo:   recordInfo,
	}

	// 获取原始链接
	originalUrl, err := h.app.Queries.GetOriginalUrl.Handle(c.Context(), q)
	if err != nil {
		return err
	}

	if originalUrl == "" {
		return error_no.ShortLinkNotFound
	}

	return c.Redirect(originalUrl)
}

// CreateShortLink 创建短链接
func (h ShortLinkApi) CreateShortLink(c *fiber.Ctx) (err error) {

	var reqParam req.ShortLinkCreateReq
	if err = c.BodyParser(&reqParam); err != nil {
		return err
	}

	cmd := command.CreateLink{
		OriginalUrl:   reqParam.OriginalUrl,
		Gid:           reqParam.Gid,
		CreateType:    reqParam.CreateType,
		ValidDateType: reqParam.ValidDateType,
		ValidDate:     reqParam.ValidDate,
		Description:   reqParam.Description,
		WithLock:      false,
	}

	if err = h.app.Commands.CreateLink.Handle(c.Context(), cmd); err != nil {
		return err
	}

	if res := cmd.ExecutionResult(); res != nil {
		response := resp.ShortLinkCreateResp{
			Gid:          res.Gid,
			OriginalUrl:  res.OriginalUrl,
			FullShortUrl: res.FullShortUrl,
		}
		return c.JSON(response)
	}
	return
}

// CreateShortLinkWithLock 加锁创建短链接
func (h ShortLinkApi) CreateShortLinkWithLock(c *fiber.Ctx) (err error) {

	var reqParam req.ShortLinkCreateReq
	if err = c.BodyParser(&reqParam); err != nil {
		return err
	}

	cmd := command.CreateLink{
		OriginalUrl:   reqParam.OriginalUrl,
		Gid:           reqParam.Gid,
		CreateType:    reqParam.CreateType,
		ValidDateType: reqParam.ValidDateType,
		ValidDate:     reqParam.ValidDate,
		Description:   reqParam.Description,
		WithLock:      true,
	}

	if err = h.app.Commands.CreateLink.Handle(c.Context(), cmd); err != nil {
		return err
	}

	if res := cmd.ExecutionResult(); res != nil {
		response := resp.ShortLinkCreateResp{
			Gid:          res.Gid,
			OriginalUrl:  res.OriginalUrl,
			FullShortUrl: res.FullShortUrl,
		}
		return c.JSON(response)
	}
	return
}

// BatchCreateShortLink 批量创建短链接
func (h ShortLinkApi) BatchCreateShortLink(c *fiber.Ctx) error {

	var reqParam req.ShortLinkBatchCreateReq
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	cmd := command.CreateLinkBatch{
		OriginUrls:    reqParam.OriginalUrls,
		Descriptions:  reqParam.Descriptions,
		Gid:           reqParam.Gid,
		CreateType:    reqParam.CreateType,
		ValidDateType: reqParam.ValidDateType,
		ValidDate:     reqParam.ValidDate,
	}

	if err := h.app.Commands.CreateLinkBatch.Handle(c.Context(), cmd); err != nil {
		return err
	}

	if res := cmd.ExecutionResult(); res != nil {
		var response resp.ShortLinkBatchCreateResp
		if err := copier.Copy(&response, &res); err != nil {
			return err
		}
		return c.JSON(response)
	}

	return nil
}

// UpdateShortLink 更新短链接
func (h ShortLinkApi) UpdateShortLink(c *fiber.Ctx) error {

	var reqParam req.ShortLinkUpdateReq
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	err := h.app.Commands.UpdateLink.Handle(c.Context(), command.UpdateLink{
		FullShortUrl:  reqParam.FullShortUrl,
		OriginalUrl:   reqParam.OriginUrl,
		OriginalGid:   reqParam.OriginGid,
		Gid:           reqParam.Gid,
		ValidDateType: reqParam.ValidDateType,
		ValidDate:     reqParam.ValidDate,
		Description:   reqParam.Description,
	})
	if err != nil {
		return err
	}

	c.Status(fiber.StatusNoContent)

	return nil
}

// PageQueryShortLink 分页查询短链接
func (h ShortLinkApi) PageQueryShortLink(c *fiber.Ctx) error {

	var reqParam req.ShortLinkPageReq
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.PageLink.Handle(c.Context(), query.PageLink{
		PageReq:  reqParam.PageReq,
		Gid:      reqParam.Gid,
		OrderTag: reqParam.OrderTag,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkPageResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(response)
}

// ListGroupShortLinkCount 查询短链接分组内数量
func (h ShortLinkApi) ListGroupShortLinkCount(c *fiber.Ctx) error {

	gidStr := c.Query("gid")
	if gidStr == "" {
		return error_no.ErrBadRequest
	}
	gidList := strings.Split(gidStr, ",")

	res, err := h.app.Queries.ListGroupCount.Handle(c.Context(), gidList)
	if err != nil {
		return err
	}

	var response resp.ShortLinkGroupCountQueryResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(response)
}
