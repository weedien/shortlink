package resp

import "time"

// ShortLinkCreateResp 短链接创建响应
type ShortLinkCreateResp struct {
	// 分组标识
	Gid string `json:"gid"`
	// 原始链接
	OriginUrl string `json:"originUrl"`
	// 短链接
	ShortUrl string `json:"shortUrl"`
}

// ShortLinkBatchCreateResp 短链接批量创建响应
type ShortLinkBatchCreateResp struct {
	// 成功数量
	SuccessCount int `json:"successCount"`
	// 批量创建返回参数
	LinkBaseInfos []ShortLinkBaseInfoDTO `json:"linkBaseInfos"`
}

// ShortLinkGroupCountQueryResp 短链接分组数量查询响应
type ShortLinkGroupCountQueryResp struct {
	// 分组标识
	Gid string `json:"gid"`
	// 短链接数量
	ShortLinkCount int `json:"shortLinkCount"`
}

// ShortLinkPageResp 短链接分页查询响应
type ShortLinkPageResp struct {
	// ID
	ID int64 `json:"id"`
	// 域名
	Domain string `json:"domain"`
	// 短链接
	ShortUri string `json:"shortUri"`
	// 完整短链接
	FullShortUrl string `json:"fullShortUrl"`
	// 原始链接
	OriginUrl string `json:"originUrl"`
	// 分组标识
	Gid string `json:"gid"`
	// 有效期类型 0：永久有效 1：自定义
	ValidDateType int `json:"validDateType"`
	// 启用标识 0：启用 1：未启用
	EnableStatus int `json:"enableStatus"`
	// 有效期
	ValidDate time.Time `json:"validDate" format:"2006-01-02 15:04:05"`
	// 创建时间
	CreateTime time.Time `json:"createTime" format:"2006-01-02 15:04:05"`
	// 描述
	Description string `json:"description"`
	// 网站标识
	Favicon string `json:"favicon"`
	// 历史PV
	TotalPv int `json:"totalPv"`
	// 今日PV
	TodayPv int `json:"todayPv"`
	// 历史UV
	TotalUv int `json:"totalUv"`
	// 今日UV
	TodayUv int `json:"todayUv"`
	// 历史UIP
	TotalUip int `json:"totalUip"`
	// 今日UIP
	TodayUip int `json:"todayUip"`
}

// ShortLinkStatsAccessRecordResp 短链接监控访问统计记录响应
type ShortLinkStatsAccessRecordResp struct {
	// 访客类型
	VisitorType string `json:"visitorType"`
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

// ShortLinkBaseInfoDTO 短链接基本信息响应
type ShortLinkBaseInfoDTO struct {
	// 描述信息
	Description string `json:"description"`
	// 原始链接
	OriginUrl string `json:"originUrl"`
	// 短链接
	ShortUrl string `json:"shortUrl"`
}
