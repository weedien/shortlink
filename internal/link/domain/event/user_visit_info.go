package event

import (
	"shortlink/internal/common/base_event"
	"time"
)

type UserVisitInfo struct {
	// 短链接
	ShortUri string `json:"shortUri"`
	// 完整短链接
	FullShortUrl string `json:"fullShortUrl"`
	// 访问用户IP
	RemoteAddr string `json:"remoteAddr"`
	// 操作系统
	OS string `json:"os"`
	// 浏览器
	Browser string `json:"browser"`
	// 操作设备
	Device string `json:"device"`
	// 网络
	Network string `json:"network"`
	// UV
	UV string `json:"uv"`
	// UV访问标识
	UVFirstFlag bool `json:"uvFirstFlag"`
	// UIP访问标识
	UipFirstFlag bool `json:"uipFirstFlag"`
	// 消息队列唯一标识
	Keys string `json:"keys"`
	// 当前时间
	CurrentDate time.Time `json:"currentDate"`
}

type UserVisitEvent struct {
	base_event.CommonEvent
	VisitInfo UserVisitInfo
}

func (e UserVisitEvent) Name() string {
	return "UserVisitEvent"
}

func NewUserVisitEvent(recordInfo UserVisitInfo) UserVisitEvent {
	return UserVisitEvent{
		CommonEvent: base_event.NewCommonEvent(),
		VisitInfo:   recordInfo,
	}
}
