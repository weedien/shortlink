package trigger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"shortlink/internal/link_stats/app"
	"shortlink/internal/link_stats/app/query"
	"shortlink/internal/link_stats/trigger/dto/req"
	"shortlink/internal/link_stats/trigger/dto/resp"
)

type ShortLinkStatsApi struct {
	app app.Application
}

func NewShortLinkStatsApi(app app.Application, router fiber.Router) {
	api := &ShortLinkStatsApi{
		app: app,
	}

	// 访问单个短链接指定时间内监控数据
	router.Get("/stats", api.GetLinkStats)
	// 访问分组短链接指定时间内监控数据
	router.Get("/stats/group", api.GroupLinkStats)
	// 访问单个短链接指定时间内访问记录监控数据
	router.Get("/stats/access-record", api.GetLinkStatsAccessRecord)
	// 访问分组短链接指定时间内访问记录监控数据
	router.Get("/stats/access-record/group", api.GroupLinkStatsAccessRecord)
}

// GetLinkStats 获取短链接统计信息
func (h ShortLinkStatsApi) GetLinkStats(c *fiber.Ctx) error {

	var reqParam req.ShortLinkStatsReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GetLinkStats.Handle(c.Context(), query.GetLinkStats{
		FullShortUrl: reqParam.FullShortUrl,
		Gid:          reqParam.Gid,
		StartDate:    reqParam.StartTime,
		EndDate:      reqParam.EndTime,
		EnableStatus: reqParam.EnableStatus,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatsResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GroupLinkStats 分组获取短链接统计信息
func (h ShortLinkStatsApi) GroupLinkStats(c *fiber.Ctx) error {

	var reqParam req.ShortLinkGroupStatsReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GroupLinkStats.Handle(c.Context(), query.GroupLinkStats{
		Gid:       reqParam.Gid,
		StartDate: reqParam.StartTime,
		EndDate:   reqParam.EndTime,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatsResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GetLinkStatsAccessRecord 获取短链接访问记录
func (h ShortLinkStatsApi) GetLinkStatsAccessRecord(c *fiber.Ctx) error {

	var reqParam req.ShortLinkStatsAccessRecordReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GetLinkStatsAccessRecord.Handle(c.Context(), query.GetLinkStatsAccessRecord{
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

	var response resp.ShortLinkStatsAccessRecordResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(res)
}

// GroupLinkStatsAccessRecord 分组获取短链接访问记录
func (h ShortLinkStatsApi) GroupLinkStatsAccessRecord(c *fiber.Ctx) error {

	var reqParam req.ShortLinkGroupStatsAccessRecordReq
	if err := c.QueryParser(&reqParam); err != nil {
		return err
	}

	res, err := h.app.Queries.GroupLinkStatsAccessRecord.Handle(c.Context(), query.GroupLinkStatsAccessRecord{
		PageReq:   reqParam.PageReq,
		Gid:       reqParam.Gid,
		StartDate: reqParam.StartTime,
		EndDate:   reqParam.EndTime,
	})
	if err != nil {
		return err
	}

	var response resp.ShortLinkStatsAccessRecordResp
	if err = copier.Copy(&response, &res); err != nil {
		return err
	}

	return c.JSON(response)
}
