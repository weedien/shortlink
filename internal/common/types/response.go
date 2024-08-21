package types

import (
	"shortlink/internal/common/error_no"
)

const SuccessCode = 10000

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func OkWithData(data interface{}) *Response {
	return &Response{
		Code: SuccessCode,
		Data: data,
	}
}

func FailWithErrorCode(errorCode error_no.ErrorCode) *Response {
	return &Response{
		Code: errorCode.Code,
		Msg:  errorCode.Desc,
		Data: nil,
	}
}

func FailWithErrorCodeAndMsg(errorCode error_no.ErrorCode, msg string) *Response {
	return &Response{
		Code: errorCode.Code,
		Msg:  msg,
		Data: nil,
	}
}

func FailWithMsg(msg string) *Response {
	return &Response{
		Code: error_no.Unknown.Code,
		Msg:  msg,
	}
}

func FailWithCodeAndMsg(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}
