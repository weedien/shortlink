package error_no

type SlugError struct {
	errorCode ErrorCode
	errorType ErrorType
	msg       string
	err       error
}

func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

//func NewServiceError(errorCode ErrorCode) SlugError {
//	return SlugError{
//		errorCode: errorCode,
//		errorType: ServiceError,
//	}
//}

func NewServiceErrorWithMsg(errorCode ErrorCode, msg string) SlugError {
	return SlugError{
		errorCode: errorCode,
		errorType: ServiceError,
		msg:       msg,
	}
}

func NewRequestError(errorCode ErrorCode) SlugError {
	return SlugError{
		errorCode: errorCode,
		errorType: RequestError,
	}
}

func NewExternalError(errorCode ErrorCode) SlugError {
	return SlugError{
		errorCode: errorCode,
		errorType: ExternalError,
	}
}

func NewExternalErrorWithMsg(errorCode ErrorCode, msg string) SlugError {
	return SlugError{
		errorCode: errorCode,
		errorType: ExternalError,
		msg:       msg,
	}
}

func (s SlugError) Error() string {
	return s.msg
}

type ErrorType string

var (
	ErrorTypeIncorrectInput   ErrorType = "incorrect-input"
	ErrorTypeAuthorization    ErrorType = "authorization"
	ErrorTypeResourceNotFound ErrorType = "resource-not-found"

	RequestError  ErrorType = "RequestError"
	ExternalError ErrorType = "ExternalError"
	ServiceError  ErrorType = "ServiceError"
)

type ErrorCode struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

var (
	//RouteNotFound = ErrorCode{Code: 10001, Desc: "请求路由不存在"}
	Unknown = ErrorCode{Code: 10015, Desc: "未知错误"}
	//TooManyRequests = ErrorCode{Code: 10016, Desc: "单位时间内请求过多，请稍后再试"}

	//OriginalUrlMisspelled    = ErrorCode{Code: 30001, Desc: "原始链接拼写错误"}
	InvalidDomain       = ErrorCode{Code: 30002, Desc: "不合法的域名"}
	LinkGenerateFailed  = ErrorCode{Code: 30003, Desc: "短链接生成失败"}
	LinkCreateFailed    = ErrorCode{Code: 30004, Desc: "短链接创建失败"}
	TooManyLinkCreate   = ErrorCode{Code: 30005, Desc: "单位时间内创建短链接过多，请稍后再试"}
	LinkDuplicateInsert = ErrorCode{Code: 30006, Desc: "短链接重复插入"}
	FaviconGetFailed    = ErrorCode{Code: 30101, Desc: "获取网站图标失败"}

	//LinkNotExists    = ErrorCode{Code: 30007, Desc: "短链接不存在"}
	LinkUpdateFailed = ErrorCode{Code: 30008, Desc: "短链接更新失败"}

	DatabaseError     = ErrorCode{Code: 40001, Desc: "数据库错误"}
	RedisConnectError = ErrorCode{Code: 40000, Desc: "Redis连接错误"}
	//RedisError        = ErrorCode{Code: 40001, Desc: "Redis错误"}
	RocketMQError = ErrorCode{Code: 40002, Desc: "RocketMQ错误"}
)
