package query

import (
	"context"
	"shortlink/internal/user/domain/user"
)

type CheckLoginCommand struct {
	Username string
	Token    string
}

type CheckLoginHandler struct {
	repo user.Repository
}

func (h CheckLoginHandler) Handle(ctx context.Context, username, token string) (bool, error) {
	return h.repo.CheckLogin(ctx, username, token)
}
