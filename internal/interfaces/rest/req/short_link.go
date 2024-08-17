package req

import "time"

// ShortLinkCreateReq 创建短链接请求
type ShortLinkCreateReq struct {
	// 域名
	Domain string `json:"domain" binding:"required"`
	// 原始链接
	OriginUrl string `json:"origin_url" binding:"required"`
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int `json:"create_type" binding:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int `json:"valid_date_type" binding:"required"`
	// 有效期
	ValidDate time.Time `json:"valid_date" binding:"required" format:"2006-01-02 15:04:05"`
	// 描述
	Description string `json:"description" binding:"required"`
}

// ShortLinkBatchCreateReq 批量创建短链接请求
type ShortLinkBatchCreateReq struct {
	// 原始链接集合
	OriginUrls []string `json:"origin_urls" binding:"required"`
	// 描述集合
	Descriptions []string `json:"descriptions" binding:"required"`
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int `json:"create_type" binding:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int `json:"valid_date_type" binding:"required"`
	// 有效期
	ValidDate time.Time `json:"valid_date" binding:"required" format:"2006-01-02 15:04:05"`
}

// ShortLinkUpdateReq 更新短链接请求
type ShortLinkUpdateReq struct {
	// 原始链接
	OriginUrl string `json:"origin_url" binding:"required"`
	// 完整短链接
	FullShortUrl string `json:"full_short_url" binding:"required"`
	// 原始分组标识
	OriginGid string `json:"origin_gid" binding:"required"`
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int `json:"valid_date_type" binding:"required"`
	// 有效期
	ValidDate time.Time `json:"valid_date" binding:"required" format:"2006-01-02 15:04:05"`
	// 描述
	Description string `json:"description" binding:"required"`
}

// ShortLinkGroupStatsAccessRecordReq 分组短链接监控访问记录请求
type ShortLinkGroupStatsAccessRecordReq struct {
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" binding:"required" format:"2006-01-02 15:04:05"`
}

// ShortLinkGroupStatsReq 分组短链接监控请求
type ShortLinkGroupStatsReq struct {
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" binding:"required" format:"2006-01-02 15:04:05"`
}

// ShortLinkPageReq 分页查询短链接请求
type ShortLinkPageReq struct {
	// 分组ID
	Gid string `json:"gid" binding:"required"`
	// 排序标识
	OrderTag string `json:"order_tag" binding:"required"`
}

// ShortLinkStatsAccessRecordReq 短链接监控访问记录请求
type ShortLinkStatsAccessRecordReq struct {
	// 完整短链接
	FullShortUrl string `json:"full_short_url" binding:"required"`
	// 分组标识
	Gid string `json:"gid" binding:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 启用标识
	EnableStatus int `json:"enable_status" binding:"required"`
}

// ShortLinkStatsReq 短链接监控请求
type ShortLinkStatsReq struct {
	// 完整短链接
	FullShortUrl string `json:"full_short_url" binding:"required"`
	// 分组标识
	Gid string `json:"gid" binding:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" binding:"required" format:"2006-01-02 15:04:05"`
	// 启用标识
	EnableStatus int `json:"enable_status" binding:"required"`
}
