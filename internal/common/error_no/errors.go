package error_no

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal error")

	ShortLinkNotExists    = errors.New("short link not exists")
	ShortLinkDisabled     = errors.New("short link disabled")
	ShortLinkExpired      = errors.New("short link expired")
	ShortLinkForbidden    = errors.New("short link forbidden")
	ShortLinkReserved     = errors.New("short link reserved")
	OriginalUrlMisspelled = errors.New("original url misspelled")
	InvalidLinkStatus     = errors.New("invalid link status")

	LockAcquireFailed = errors.New("lock acquire failed")
	LockReleaseFailed = errors.New("lock release failed")

	RouteNotFound   = errors.New("route not found")
	TooManyRequests = errors.New("too many requests")

	UserNotExist                    = errors.New("user not exist")
	UserExist                       = errors.New("user exist")
	UserForbidden                   = errors.New("user forbidden")
	UserUnlogged                    = errors.New("user unlogged")
	InvalidTokenOrUnloggedLoginUser = errors.New("invalid token or unlogged user")

	RedisError       = errors.New("redis error")
	RedisKeyNotExist = errors.New("redis key not exist")
	RedisKeyExpired  = errors.New("redis key expired")
)

func NewServiceError(err error) SlugError {
	return SlugError{
		err:       err,
		errorType: ServiceError,
	}
}
