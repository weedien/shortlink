package command

import (
	"context"
	"shortlink/internal/common/error_no"
	"shortlink/internal/user/domain/user"
)

type UserLogoutCommand struct {
	Username string
	Token    string
}

type UserLogoutHandler struct {
	repo user.Repository
}

func (h UserLogoutHandler) Handle(ctx context.Context, username, token string) (err error) {
	var login bool
	if login, err = h.repo.CheckLogin(ctx, username, token); err != nil {
		return err
	}
	if !login {
		return error_no.InvalidTokenOrUnloggedLoginUser
	}
	return h.repo.InvalidateToken(ctx, username, token)
}
