package valobj

import (
	"shortlink/internal/common/types"
	"time"
)

type ShortLinkCreateVo struct {
	FullShortUrl string
	OriginalUrl  string
	Gid          string
}

type ShortLinkCreateBatchVo struct {
	// 成功数量
	SuccessCount int
	// 批量创建返回参数
	LinkInfos []ShortLinkCreateVo
}

// ShortLinkQueryVo 这两个字段能确定一个唯一的短链
type ShortLinkQueryVo struct {
	Gid          string
	FullShortUrl string
}

type ShortLinkUpdateVo struct {
	// 完整短链接
	FullShortLink string
	// 原始链接
	OriginalUrl string
	// 原始分组ID
	OriginalGid string
	// 分组ID
	Gid string
	// 有效期类型 0:永久有效 1:自定义有效期
	ValidDateType int
	// 有效期
	ValidDate time.Time
	// 描述
	Description string
}

func (v ShortLinkUpdateVo) GidChanged() bool {
	return v.OriginalGid != v.Gid
}

type ShortLinkPageQueryVO struct {
	// 分页请求
	types.PageReq
	// 分组ID
	Gid string
	// 排序标识
	OrderTag string
}

type ShortLinkPageRespVO struct {
	// 分页数据
	types.PageResp[ShortLinkQueryVo]
}

type TodayStatsVo struct {
	// 今日PV
	TodayPv int
	// 今日UV
	TodayUv int
	// 今日UIP
	TodayUip int
}

type ShortLinkStatsRecordVo struct {
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

type RecordLinkVisitInfoVo struct {
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
	UIPFirstFlag bool `json:"uipFirstFlag"`
	// 消息队列唯一标识
	Keys string `json:"keys"`
	// 当前时间
	CurrentDate time.Time `json:"currentDate"`
}
