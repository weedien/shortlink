package resp

import (
	"shortlink/internal/common/types"
	"time"
)

// ShortLinkStatsResp 短链接监控统计响应
type ShortLinkStatsResp struct {
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
	// 小时访问统计
	Hourly []ShortLinkStatsAccessBaseDTO `json:"hourly"`
	// 日访问统计
	Daily []ShortLinkStatsAccessBaseDTO `json:"daily"`
	// 周访问统计
	Weekly []ShortLinkStatsAccessBaseDTO `json:"weekly"`
	// 地区访问统计（仅国内）
	LocationCnStats []ShortLinkStatsLocationDTO `json:"locationCnStats"`
	// 高频访问IP统计
	TopIpStats []ShortLinkStatsTopIpDTO `json:"topIpStats"`
	// 浏览器统计
	BrowserStats []ShortLinkStatsBrowserDTO `json:"browserStats"`
	// 操作系统统计
	OsStats []ShortLinkStatsOsDTO `json:"osStats"`
	// 访客类型统计
	VisitorTypeStats []ShortLinkStatsUvDTO `json:"visitorTypeStats"`
	// 设备统计
	DeviceStats []ShortLinkStatsDeviceDTO `json:"deviceStats"`
	// 网络统计
	NetworkStats []ShortLinkStatsNetworkDTO `json:"networkStats"`
}

// ShortLinkStatsAccessBaseDTO 短链接监控访问统计基础响应
type ShortLinkStatsAccessBaseDTO struct {
	// 日期
	Date time.Time `json:"date" format:"2006-01-02"`
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
}

// ShortLinkStatsBrowserDTO 浏览器统计响应
type ShortLinkStatsBrowserDTO struct {
	// 统计
	Count int `json:"count"`
	// 浏览器
	Browser string `json:"browser"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsDeviceDTO 设备统计响应
type ShortLinkStatsDeviceDTO struct {
	// 统计
	Count int `json:"count"`
	// 设备
	Device string `json:"device"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsLocationDTO 地区统计响应
type ShortLinkStatsLocationDTO struct {
	// 统计
	Count int `json:"count"`
	// 地区
	Location string `json:"location"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsNetworkDTO 网络统计响应
type ShortLinkStatsNetworkDTO struct {
	// 统计
	Count int `json:"count"`
	// 网络
	Network string `json:"network"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsOsDTO 操作系统统计响应
type ShortLinkStatsOsDTO struct {
	// 统计
	Count int `json:"count"`
	// 操作系统
	Os string `json:"os"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsTopIpDTO 短链接高频访问IP统计响应
type ShortLinkStatsTopIpDTO struct {
	// 统计
	Count int `json:"count"`
	// IP
	Ip string `json:"ip"`
}

// ShortLinkStatsUvDTO 短链接监控访问统计UV响应
type ShortLinkStatsUvDTO struct {
	// 统计
	Count int `json:"count"`
	// 访客类型
	VisitorType string `json:"visitorType"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsAccessRecordResp 短链接监控访问统计记录响应
type ShortLinkStatsAccessRecordResp types.PageResp[ShortLinkStatsAccessRecordDTO]

type ShortLinkStatsAccessRecordDTO struct {
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
