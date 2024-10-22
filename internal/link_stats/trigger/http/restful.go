package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"shortlink/internal/link_stats/app"
	"shortlink/internal/link_stats/app/query"
	"shortlink/internal/link_stats/trigger/http/dto/req"
	"shortlink/internal/link_stats/trigger/http/dto/resp"
)

type ShortLinkStatApi struct {
	app app.Application
}

func NewShortLinkStatApi(app app.Application, router fiber.Router) {
	api := &ShortLinkStatApi{
		app: app,
	}

	// 访问单个短链接指定时间内监控数据
	router.Get("/stats", api.GetLinkStat)
	// 访问分组短链接指定时间内监控数据
	router.Get("/stats/group", api.GroupLinkStat)
	// 访问单个短链接指定时间内访问记录监控数据
	router.Get("/stats/access-record", api.GetLinkStatAccessRecord)
	// 访问分组短链接指定时间内访问记录监控数据
	router.Get("/stats/access-record/group", api.GroupLinkStatAccessRecord)
}

// GetLinkStat 获取短链接统计信息
func (h ShortLinkStatApi) GetLinkStat(c *fiber.Ctx) error {

	var reqParam req.ShortLinkStatReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GetLinkStat.Handle(c.Context(), query.GetLinkStat{
		FullShortUrl: reqParam.FullShortUrl,
		Gid:          reqParam.Gid,
		StartDate:    reqParam.StartTime,
		EndDate:      reqParam.EndTime,
		EnableStatus: reqParam.EnableStatus,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GroupLinkStat 分组获取短链接统计信息
func (h ShortLinkStatApi) GroupLinkStat(c *fiber.Ctx) error {

	var reqParam req.ShortLinkGroupStatReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GroupLinkStat.Handle(c.Context(), query.GroupLinkStat{
		Gid:       reqParam.Gid,
		StartDate: reqParam.StartTime,
		EndDate:   reqParam.EndTime,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GetLinkStatAccessRecord 获取短链接访问记录
func (h ShortLinkStatApi) GetLinkStatAccessRecord(c *fiber.Ctx) error {

	var reqParam req.ShortLinkStatAccessRecordReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GetLinkStatAccessRecord.Handle(c.Context(), query.GetLinkStatAccessRecord{
		PageReq:      reqParam.PageReq,
		FullShortUrl: reqParam.FullShortUrl,
		Gid:          reqParam.Gid,
		StartDate:    reqParam.StartTime,
		EndDate:      reqParam.EndTime,
		EnableStatus: reqParam.EnableStatus,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatAccessRecordResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GroupLinkStatAccessRecord 分组获取短链接访问记录
func (h ShortLinkStatApi) GroupLinkStatAccessRecord(c *fiber.Ctx) error {

	var reqParam req.ShortLinkGroupStatAccessRecordReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GroupLinkStatAccessRecord.Handle(c.Context(), query.GroupLinkStatAccessRecord{
		PageReq:   reqParam.PageReq,
		Gid:       reqParam.Gid,
		StartDate: reqParam.StartTime,
		EndDate:   reqParam.EndTime,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatAccessRecordResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(response)
}
