package http

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/internal/user/app/group"
	"shortlink/internal/user/trigger/http/dto/req"
	"shortlink/internal/user/trigger/http/dto/resp"
)

type GroupApi struct {
	app group.Application
}

func NewGroupApi(app group.Application, router fiber.Router) {
	api := &GroupApi{app: app}

	groupRouter := router.Group("/group")
	groupRouter.Get("", api.List)
	groupRouter.Post("", api.Create)
	groupRouter.Post("/sort", api.Sort)
	groupRouter.Put("", api.Update)
	groupRouter.Delete("/:gid", api.Delete)
}

// Create 新增短链接分组
func (h GroupApi) Create(c *fiber.Ctx) error {
	reqParam := req.ShortLinkGroupSaveReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}
	return nil
}

// List 查询短链接分组集合
func (h GroupApi) List(c *fiber.Ctx) error {
	response := resp.ShortLinkGroupResp{}

	return c.JSON(response)
}

// Update 修改短链接分组名称
func (h GroupApi) Update(c *fiber.Ctx) error {
	reqParam := req.ShortLinkGroupUpdateReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}
	return nil
}

// Delete 删除短链接分组
func (h GroupApi) Delete(c *fiber.Ctx) error {
	_ = c.Params("gid")

	return nil
}

// Sort 排序短链接分组
func (h GroupApi) Sort(c *fiber.Ctx) error {
	reqParam := req.ShortLinkGroupSortReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	return nil
}
