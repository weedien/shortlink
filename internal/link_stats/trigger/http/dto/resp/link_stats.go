package resp

import (
	"shortlink/internal/common/types"
	"time"
)

// ShortLinkStatResp 短链接监控统计响应
type ShortLinkStatResp struct {
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
	// 小时访问统计
	Hourly []ShortLinkStatAccessBaseDTO `json:"hourly"`
	// 日访问统计
	Daily []ShortLinkStatAccessBaseDTO `json:"daily"`
	// 周访问统计
	Weekly []ShortLinkStatAccessBaseDTO `json:"weekly"`
	// 地区访问统计（仅国内）
	LocationCnStat []ShortLinkStatLocationDTO `json:"locationCnStat"`
	// 高频访问IP统计
	TopIpStat []ShortLinkStatTopIpDTO `json:"topIpStat"`
	// 浏览器统计
	BrowserStat []ShortLinkStatBrowserDTO `json:"browserStat"`
	// 操作系统统计
	OsStat []ShortLinkStatOsDTO `json:"osStat"`
	// 访客类型统计
	VisitorTypeStat []ShortLinkStatUvDTO `json:"visitorTypeStat"`
	// 设备统计
	DeviceStat []ShortLinkStatDeviceDTO `json:"deviceStat"`
	// 网络统计
	NetworkStat []ShortLinkStatNetworkDTO `json:"networkStat"`
}

// ShortLinkStatAccessBaseDTO 短链接监控访问统计基础响应
type ShortLinkStatAccessBaseDTO struct {
	// 日期
	Date time.Time `json:"date" format:"2006-01-02"`
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
}

// ShortLinkStatBrowserDTO 浏览器统计响应
type ShortLinkStatBrowserDTO struct {
	// 统计
	Count int `json:"count"`
	// 浏览器
	Browser string `json:"browser"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatDeviceDTO 设备统计响应
type ShortLinkStatDeviceDTO struct {
	// 统计
	Count int `json:"count"`
	// 设备
	Device string `json:"device"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatLocationDTO 地区统计响应
type ShortLinkStatLocationDTO struct {
	// 统计
	Count int `json:"count"`
	// 地区
	Location string `json:"location"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatNetworkDTO 网络统计响应
type ShortLinkStatNetworkDTO struct {
	// 统计
	Count int `json:"count"`
	// 网络
	Network string `json:"network"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatOsDTO 操作系统统计响应
type ShortLinkStatOsDTO struct {
	// 统计
	Count int `json:"count"`
	// 操作系统
	Os string `json:"os"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatTopIpDTO 短链接高频访问IP统计响应
type ShortLinkStatTopIpDTO struct {
	// 统计
	Count int `json:"count"`
	// IP
	Ip string `json:"ip"`
}

// ShortLinkStatUvDTO 短链接监控访问统计UV响应
type ShortLinkStatUvDTO struct {
	// 统计
	Count int `json:"count"`
	// 访客类型
	VisitorType string `json:"visitorType"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatAccessRecordResp 短链接监控访问统计记录响应
type ShortLinkStatAccessRecordResp types.PageResp[ShortLinkStatAccessRecordDTO]

type ShortLinkStatAccessRecordDTO struct {
	// 访客类型
	UvType string `json:"UvType"`
	// 浏览器
	Browser string `json:"browser"`
	// 操作系统
	Os string `json:"os"`
	// IP
	Ip string `json:"ip"`
	// 访问网络
	Network string `json:"network"`
	// 访问设备
	Device string `json:"device"`
	// 地区
	Location string `json:"location"`
	// 用户信息
	UserAgent string `json:"userAgent"`
	// 访问时间
	AccessTime time.Time `json:"accessTime" format:"2006-01-02 15:04:05"`
}
