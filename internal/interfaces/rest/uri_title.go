package rest

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/common/toolkit"
)

type UriTitleApi struct {
}

func NewUriTitleApi(router fiber.Router) {
	server := &UriTitleApi{}

	router.Get("/get-title", server.GetTitleByUrl)
}

func (h UriTitleApi) GetTitleByUrl(c *fiber.Ctx) (err error) {
	url := c.Params("url")

	var title string
	if title, err = toolkit.GetTitleByUrl(url); err != nil {
		return err
	}
	err = c.JSON(fiber.Map{
		"title": title,
	})
	return
}
