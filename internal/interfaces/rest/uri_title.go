package rest

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/pkg/toolkit"
	"shortlink/pkg/types"
)

type UriTitleHandler struct {
}

func NewUriTitleServer(router fiber.Router) {
	server := &UriTitleHandler{}

	router.Get("/get-title", server.GetTitleByUrl)
}

func (h UriTitleHandler) GetTitleByUrl(c *fiber.Ctx) (err error) {
	url := c.Params("url")
	title, err := toolkit.GetTitleByUrl(url)
	err = c.JSON(types.OkWithData(title))
	return
}
