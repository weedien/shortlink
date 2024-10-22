package query

import "time"

// LinkStatAccessRecord 短链接监控访问记录
type LinkStatAccessRecord struct {
	// 访客类型
	UvType string
	// 浏览器
	Browser string
	// 操作系统
	Os string
	// ip
	Ip string
	// 访问网络
	Network string
	// 访问设备
	Device string
	// 地区
	Locale string
	// 用户信息
	User string
	// 访问时间
	AccessTime time.Time
}

// LinkStat 短链接监控统计
type LinkStat struct {
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
	// 小时访问统计
	Hourly []int `json:"hourly"`
	// 日访问统计
	Daily []LinkStatAccessDaily `json:"daily"`
	// 周访问统计
	Weekly []int `json:"weekly"`
	// 地区访问统计（仅国内）
	LocationCnStat []LinkStatLocale `json:"locationCnStat"`
	// 高频访问IP统计
	TopIpStat []LinkStatTopIp `json:"topIpStat"`
	// 浏览器统计
	BrowserStat []LinkStatBrowser `json:"browserStat"`
	// 操作系统统计
	OsStat []LinkStatOs `json:"osStat"`
	// 访客类型统计
	VisitorTypeStat []LinkStatUv `json:"visitorTypeStat"`
	// 设备统计
	DeviceStat []LinkStatDevice `json:"deviceStat"`
	// 网络统计
	NetworkStat []LinkStatNetwork `json:"networkStat"`
}

// LinkStatAccessDaily 短链接监控访问统计基础响应
type LinkStatAccessDaily struct {
	// 日期
	Date string `json:"date"`
	// PV
	Pv int `json:"pv"`
	// UV
	Uv int `json:"uv"`
	// UIP
	Uip int `json:"uip"`
}

// LinkStatBrowser 浏览器统计响应
type LinkStatBrowser struct {
	// 统计
	Cnt int `json:"count"`
	// 浏览器
	Browser string `json:"browser"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// LinkStatDevice 设备统计响应
type LinkStatDevice struct {
	// 统计
	Cnt int `json:"count"`
	// 设备
	Device string `json:"device"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// LinkStatLocale 地区统计响应
type LinkStatLocale struct {
	// 统计
	Cnt int `json:"count"`
	// 地区
	Locale string `json:"locale"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// LinkStatNetwork 网络统计响应
type LinkStatNetwork struct {
	// 统计
	Cnt int `json:"count"`
	// 网络
	Network string `json:"network"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// LinkStatOs 操作系统统计响应
type LinkStatOs struct {
	// 统计
	Cnt int `json:"count"`
	// 操作系统
	Os string `json:"os"`
	// 占比
	Ratio float64 `json:"ratio"`
}

// LinkStatTopIp 短链接高频访问IP统计响应
type LinkStatTopIp struct {
	// 统计
	Cnt int `json:"count"`
	// IP
	Ip string `json:"ip"`
}

// LinkStatUv 短链接监控访问统计UV响应
type LinkStatUv struct {
	// 统计
	Cnt int `json:"count"`
	// 访客类型
	VisitorType string `json:"visitorType"`
	// 占比
	Ratio float64 `json:"ratio"`
}
