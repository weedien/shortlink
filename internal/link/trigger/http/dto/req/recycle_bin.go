package req

import (
	"shortlink/internal/common/types"
)

// RecycleBinSaveReq 回收站保存请求
type RecycleBinSaveReq struct {
	// 分组标识
	Gid string `json:"gid" binding:"required"`
	// 短链接
	ShortUri string `json:"shortUri" binding:"required"`
}

// RecycleBinRecoverReq 回收站恢复请求
type RecycleBinRecoverReq struct {
	// 分组标识
	Gid string `json:"gid" binding:"required"`
	// 短链接
	ShortUri string `json:"shortUri" binding:"required"`
}

// RecycleBinDeleteReq 回收站删除请求
type RecycleBinDeleteReq struct {
	// 分组标识
	Gid string `json:"gid" binding:"required"`
	// 短链接
	ShortUri string `json:"shortUri" binding:"required"`
}

// RecycleBinPageReq 分页查询回收站请求
type RecycleBinPageReq struct {
	// 分页参数
	types.PageReq
	// 分组标识
	Gids []string `json:"gids" binding:"required"`
}
