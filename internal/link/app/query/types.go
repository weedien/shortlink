package query

import "time"

type GroupLinkCount struct {
	// 分组ID
	Gid string
	// 短链接数量
	Count int
}

type Link struct {
	ID            int
	Domain        string
	ShortUri      string
	FullShortUrl  string
	OriginalUrl   string
	ClickNum      int
	Gid           string
	EnableStatus  int
	CreateType    int
	ValidDateType int
	ValidDate     time.Time
	Desc          string
	Favicon       string
	TotalPv       int
	TotalUv       int
	TotalUip      int
	TodayPv       int
	TodayUv       int
	TodayUip      int
}

type LinkQueryDTO struct {
	ID            int
	Domain        string
	ShortUri      string
	FullShortUrl  string
	OriginalUrl   string
	ClickNum      int
	Gid           string
	EnableStatus  int
	CreateType    int
	ValidDateType int
	ValidDate     time.Time
	Desc          string
	CreateTime    time.Time
	Favicon       string
	TotalPv       int
	TotalUv       int
	TotalUip      int
	TodayPv       int
	TodayUv       int
	TodayUip      int
}
