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

// ShortLinkStat 短链接监控统计
type ShortLinkStat struct {
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
	// 小时访问统计
	HourStat []ShortLinkStatAccessBaseVO `json:"hourly"`
	// 日访问统计
	Daily []ShortLinkStatAccessBaseVO `json:"daily"`
	// 周访问统计
	Weekly []ShortLinkStatAccessBaseVO `json:"weekly"`
	// 地区访问统计（仅国内）
	LocationCnStat []ShortLinkStatLocationVO `json:"locationCnStat"`
	// 高频访问IP统计
	TopIpStat []ShortLinkStatTopIpVO `json:"topIpStat"`
	// 浏览器统计
	BrowserStat []ShortLinkStatBrowserVO `json:"browserStat"`
	// 操作系统统计
	OsStat []ShortLinkStatOsVO `json:"osStat"`
	// 访客类型统计
	VisitorTypeStat []ShortLinkStatUvVO `json:"visitorTypeStat"`
	// 设备统计
	DeviceStat []ShortLinkStatDeviceVO `json:"deviceStat"`
	// 网络统计
	NetworkStat []ShortLinkStatNetworkVO `json:"networkStat"`
}

// ShortLinkStatAccessBaseVO 短链接监控访问统计基础响应
type ShortLinkStatAccessBaseVO struct {
	// 日期
	Date time.Time `json:"date" format:"2006-01-02"`
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
}

// ShortLinkStatBrowserVO 浏览器统计响应
type ShortLinkStatBrowserVO struct {
	// 统计
	Count int `json:"count"`
	// 浏览器
	Browser string `json:"browser"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatDeviceVO 设备统计响应
type ShortLinkStatDeviceVO struct {
	// 统计
	Count int `json:"count"`
	// 设备
	Device string `json:"device"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatLocationVO 地区统计响应
type ShortLinkStatLocationVO struct {
	// 统计
	Count int `json:"count"`
	// 地区
	Location string `json:"location"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatNetworkVO 网络统计响应
type ShortLinkStatNetworkVO struct {
	// 统计
	Count int `json:"count"`
	// 网络
	Network string `json:"network"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatOsVO 操作系统统计响应
type ShortLinkStatOsVO struct {
	// 统计
	Count int `json:"count"`
	// 操作系统
	Os string `json:"os"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// ShortLinkStatTopIpVO 短链接高频访问IP统计响应
type ShortLinkStatTopIpVO struct {
	// 统计
	Count int `json:"count"`
	// IP
	Ip string `json:"ip"`
}

// ShortLinkStatUvVO 短链接监控访问统计UV响应
type ShortLinkStatUvVO struct {
	// 统计
	Count int `json:"count"`
	// 访客类型
	VisitorType string `json:"visitorType"`
	// 占比
	Ratio float64 `json:"ratio"`
}
