package httperr

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"shortlink/common/error_no"
)

type Response struct {
	Status string
	Msg    string `json:"msg"`
}

func RespondWithError(c *fiber.Ctx, err error) error {

	// 在这里去处理一些很特殊的error，比如能和http状态码直接对应
	if errors.Is(err, error_no.TooManyRequests) {
		// 请求频率过高
		r := Response{
			Status: "Too many requests",
			Msg:    err.Error(),
		}
		return c.Status(fiber.StatusTooManyRequests).JSON(r)
	}
	if errors.Is(err, error_no.RouteNotFound) {
		// 请求频率过高
		r := Response{
			Status: "Route not found",
			Msg:    err.Error(),
		}
		return c.Status(fiber.StatusNotFound).JSON(r)
	}

	var slugError error_no.SlugError
	if ok := errors.As(err, &slugError); !ok {
		// 未定义的内部异常
		r := Response{
			Status: "Internal server error",
			Msg:    err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	switch slugError.ErrorType() {
	case error_no.ErrorTypeAuthorization:
		// 未授权
		r := Response{
			Status: "Unauthorised",
			Msg:    slugError.Error(),
		}
		return c.Status(fiber.StatusUnauthorized).JSON(r)
	case error_no.ErrorTypeIncorrectInput:
		// 请求参数错误
		r := Response{
			Status: "Bad request",
			Msg:    slugError.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(r)
	case error_no.ErrorTypeResourceNotFound:
		// 资源未找到，返回 notfound 页面
		return c.Status(fiber.StatusNotFound).SendFile("resources/notfound.html")
	default:
		// 未定义的内部异常
		r := Response{
			Status: "Internal server error",
			Msg:    slugError.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}
}

// ErrorHandler 全局错误处理（fiber专用）
func ErrorHandler(c *fiber.Ctx, err error) error {
	return RespondWithError(c, err)
}
