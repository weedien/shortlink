package error_no

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal error")

	ShortLinkNotFound     = errors.New("short link not found")
	ShortLinkExpired      = errors.New("short link expired")
	OriginalUrlMisspelled = errors.New("original url misspelled")

	LockAcquireFailed = errors.New("lock acquire failed")
	LockReleaseFailed = errors.New("lock release failed")

	RouteNotFound   = errors.New("route not found")
	TooManyRequests = errors.New("too many requests")

	UserNotExist                    = errors.New("user not exist")
	UserExist                       = errors.New("user exist")
	UserForbidden                   = errors.New("user forbidden")
	UserUnlogged                    = errors.New("user unlogged")
	InvalidTokenOrUnloggedLoginUser = errors.New("invalid token or unlogged user")
)

func NewServiceError(err error) SlugError {
	return SlugError{
		err:       err,
		errorType: ServiceError,
	}
}
