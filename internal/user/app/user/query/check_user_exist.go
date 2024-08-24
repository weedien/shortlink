package query

import "shortlink/internal/user/domain"

type CheckUserExistHandler struct {
	repo domain.UserRepository
}
