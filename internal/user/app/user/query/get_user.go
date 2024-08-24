package query

import (
	"context"
	"shortlink/internal/common/error_no"
	"shortlink/internal/user/domain"
)

type GetUserHandler struct {
	repo domain.UserRepository
}

func NewGetUserHandler(repo domain.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h GetUserHandler) Handle(ctx context.Context, username string) (res *domain.User, err error) {
	if res, err = h.repo.GetUser(ctx, username); err != nil {
		return
	}
	if res == nil {
		return nil, error_no.UserNotExist
	}
	return
}
