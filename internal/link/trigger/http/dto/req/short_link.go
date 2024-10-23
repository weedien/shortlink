package req

import (
	"shortlink/internal/common/types"
	"time"
)

// LinkCreateReq 创建短链接请求
type LinkCreateReq struct {
	// 域名
	//Domain string `json:"domain" validate:"required"`
	// 原始链接
	OriginalUrl string `json:"original_url" validate:"required"`
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int `json:"create_type" validate:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidType int `json:"valid_date_type" validate:"required"`
	// 有效期 - 开始时间
	StartDate time.Time `json:"start_date" validate:"required" format:"2006-01-02 15:04:05"`
	// 有效期 - 结束时间
	EndDate time.Time `json:"end_date" validate:"required" format:"2006-01-02 15:04:05"`
	// 描述
	Desc string `json:"desc" validate:"required"`
}

// LinkBatchCreateReq 批量创建短链接请求
type LinkBatchCreateReq struct {
	// 原始链接集合
	OriginalUrls []string `json:"original_urls" validate:"required"`
	// 描述集合
	Descs []string `json:"descs" validate:"required"`
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 创建类型 0:接口创建 1:控制台创建
	CreateType int `json:"create_type" validate:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidType int `json:"valid_date_type" validate:"required"`
	// 有效期 - 开始时间
	StartDate time.Time `json:"start_date" validate:"required" format:"2006-01-02 15:04:05"`
	// 有效期 - 结束时间
	EndDate time.Time `json:"end_date" validate:"required" format:"2006-01-02 15:04:05"`
}

// LinkUpdateReq 更新短链接请求
type LinkUpdateReq struct {
	// 原始链接
	OriginalUrl string `json:"original_url" validate:"required"`
	// 短链接
	ShortUri string `json:"short_uri" validate:"required"`
	// 原始分组标识
	OriginalGid string `json:"original_gid" validate:"required"`
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 状态
	Status string `json:"status" validate:"required"`
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidType int `json:"valid_date_type" validate:"required"`
	// 有效期 - 开始时间
	StartDate time.Time `json:"start_date" validate:"required" format:"2006-01-02 15:04:05"`
	// 有效期 - 结束时间
	EndDate time.Time `json:"end_date" validate:"required" format:"2006-01-02 15:04:05"`
	// 描述
	Desc string `json:"desc" validate:"required"`
}

// LinkGroupStatAccessRecordReq 分组短链接监控访问记录请求
type LinkGroupStatAccessRecordReq struct {
	// 分页参数
	types.PageReq `json:",inline" validate:"required"`
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" validate:"required" format:"2006-01-02 15:04:05"`
}

// LinkGroupStatReq 分组短链接监控请求
type LinkGroupStatReq struct {
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" validate:"required" format:"2006-01-02 15:04:05"`
}

// LinkPageReq 分页查询短链接请求
type LinkPageReq struct {
	// 分页参数
	types.PageReq `json:",inline" validate:"required"`
	// 分组ID
	Gid string `json:"gid" validate:"required"`
	// 排序标识
	OrderTag string `json:"order_tag" validate:"required"`
}

// LinkStatsAccessRecordReq 短链接监控访问记录请求
type LinkStatsAccessRecordReq struct {
	// 分页参数
	types.PageReq `json:",inline" validate:"required"`
	// 短链接
	ShortUri string `json:"short_uri" validate:"required"`
	// 分组标识
	Gid string `json:"gid" validate:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 启用标识
	EnableStatus int `json:"enable_status" validate:"required"`
}

// LinkStatsReq 短链接监控请求
type LinkStatsReq struct {
	// 完整短链接
	ShortUri string `json:"short_uri" validate:"required"`
	// 分组标识
	Gid string `json:"gid" validate:"required"`
	// 开始时间
	StartTime time.Time `json:"start_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 结束时间
	EndTime time.Time `json:"end_time" validate:"required" format:"2006-01-02 15:04:05"`
	// 启用标识
	EnableStatus int `json:"enable_status" validate:"required"`
}
