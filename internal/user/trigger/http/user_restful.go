package http

import (
	"github.com/gofiber/fiber/v2"
	"shortlink/internal/user/app/user"
	"shortlink/internal/user/trigger/http/dto/req"
	"shortlink/internal/user/trigger/http/dto/resp"
)

type UserApi struct {
	app user.UserApplication
}

func NewUserApi(app user.UserApplication, router fiber.Router) {
	api := &UserApi{app: app}

	userRouter := router.Group("/user")
	userRouter.Get("/:username", api.GetUserByUsername)
	userRouter.Get("/actual/:username", api.GetUserByUsernameWithoutMask)
	userRouter.Get("/exist", api.CheckUserExist)
	userRouter.Get("/check-login", api.CheckLogin)
	userRouter.Post("", api.Register)
	userRouter.Post("/login", api.Login)
	userRouter.Post("/logout", api.Logout)
	userRouter.Put("", api.Update)
	userRouter.Delete("/:username", api.Delete)
}

// GetUserByUsername 根据用户名查询用户信息
func (h UserApi) GetUserByUsername(c *fiber.Ctx) error {
	_ = c.Params("username")
	response := resp.UserResp{}

	return c.JSON(response)
}

// GetUserByUsernameWithoutMask 根据用户名查询无脱敏用户信息
func (h UserApi) GetUserByUsernameWithoutMask(c *fiber.Ctx) error {
	_ = c.Params("username")
	response := resp.UserActualResp{}

	return c.JSON(response)
}

// CheckUserExist 查询用户是否存在
func (h UserApi) CheckUserExist(c *fiber.Ctx) error {

	_ = c.Query("username")

	return c.JSON(false)
}

// Register 用户注册
func (h UserApi) Register(c *fiber.Ctx) error {
	reqParam := req.UserRegisterReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	return nil
}

// Update 修改用户信息
func (h UserApi) Update(c *fiber.Ctx) error {
	reqParam := req.UserUpdateReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}

	return nil
}

// Login 用户登录
func (h UserApi) Login(c *fiber.Ctx) error {
	reqParam := req.UserLoginReq{}
	if err := c.BodyParser(&reqParam); err != nil {
		return err
	}
	response := resp.UserLoginResp{}

	return c.JSON(response)
}

// CheckLogin 检查用户是否登录
func (h UserApi) CheckLogin(c *fiber.Ctx) error {
	_ = c.Query("username")

	return c.JSON(false)
}

// Logout 用户登出
func (h UserApi) Logout(c *fiber.Ctx) error {
	_ = c.Query("username")

	return nil
}

// Delete 删除用户
func (h UserApi) Delete(c *fiber.Ctx) error {
	_ = c.Params("username")

	return nil
}
