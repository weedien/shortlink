package command

import (
	"context"
	"shortlink/internal/common/error_no"
	"shortlink/internal/user/domain/user"
)

type UpdateUserCommand struct {
	Username string
	Password string
	RealName string
	Email    string
	Phone    string
}

type UpdateUserHandler struct {
	repo user.Repository
}

func (h UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) error {
	currentUsername := ctx.Value("username").(string)
	if currentUsername != cmd.Username {
		return error_no.UserForbidden
	}

	user := user.NewUser(cmd.Username, cmd.Password, cmd.RealName, cmd.Email, cmd.Phone)
	return h.repo.UpdateUser(ctx, &user)
}
