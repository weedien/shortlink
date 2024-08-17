package valobj

import "time"

type ShortLinkSimpleVO struct {
	// 完整短链接
	FullShortUrl string
	// 分组ID
	Gid string
	// 开始日期
	StartDate time.Time
	// 结束日期
	EndDate time.Time
	// 启用标识
	EnableStatus int
}

// ShortLinkStats 短链接监控统计
type ShortLinkStats struct {
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
	// 小时访问统计
	HourStats []ShortLinkStatsAccessBaseVO `json:"hourly"`
	// 日访问统计
	Daily []ShortLinkStatsAccessBaseVO `json:"daily"`
	// 周访问统计
	Weekly []ShortLinkStatsAccessBaseVO `json:"weekly"`
	// 地区访问统计（仅国内）
	LocationCnStats []ShortLinkStatsLocationVO `json:"locationCnStats"`
	// 高频访问IP统计
	TopIpStats []ShortLinkStatsTopIpVO `json:"topIpStats"`
	// 浏览器统计
	BrowserStats []ShortLinkStatsBrowserVO `json:"browserStats"`
	// 操作系统统计
	OsStats []ShortLinkStatsOsVO `json:"osStats"`
	// 访客类型统计
	VisitorTypeStats []ShortLinkStatsUvVO `json:"visitorTypeStats"`
	// 设备统计
	DeviceStats []ShortLinkStatsDeviceVO `json:"deviceStats"`
	// 网络统计
	NetworkStats []ShortLinkStatsNetworkVO `json:"networkStats"`
}

// ShortLinkStatsAccessBaseVO 短链接监控访问统计基础响应
type ShortLinkStatsAccessBaseVO struct {
	// 日期
	Date time.Time `json:"date" format:"2006-01-02"`
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
}

// ShortLinkStatsBrowserVO 浏览器统计响应
type ShortLinkStatsBrowserVO struct {
	// 统计
	Count int `json:"count"`
	// 浏览器
	Browser string `json:"browser"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsDeviceVO 设备统计响应
type ShortLinkStatsDeviceVO struct {
	// 统计
	Count int `json:"count"`
	// 设备
	Device string `json:"device"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsLocationVO 地区统计响应
type ShortLinkStatsLocationVO struct {
	// 统计
	Count int `json:"count"`
	// 地区
	Location string `json:"location"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsNetworkVO 网络统计响应
type ShortLinkStatsNetworkVO struct {
	// 统计
	Count int `json:"count"`
	// 网络
	Network string `json:"network"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsOsVO 操作系统统计响应
type ShortLinkStatsOsVO struct {
	// 统计
	Count int `json:"count"`
	// 操作系统
	Os string `json:"os"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatsTopIpVO 短链接高频访问IP统计响应
type ShortLinkStatsTopIpVO struct {
	// 统计
	Count int `json:"count"`
	// IP
	Ip string `json:"ip"`
}

// ShortLinkStatsUvVO 短链接监控访问统计UV响应
type ShortLinkStatsUvVO struct {
	// 统计
	Count int `json:"count"`
	// 访客类型
	VisitorType string `json:"visitorType"`
	// 占比
	Ratio float64 `json:"ratio"`
}
