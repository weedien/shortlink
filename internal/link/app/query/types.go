package query

import "time"

type GroupLinkCount struct {
	// 分组ID
	Gid string
	// 短链接数量
	Count int
}

type Link struct {
	ID           int
	Domain       string
	ShortUri     string
	FullShortUrl string
	OriginalUrl  string
	Gid          string
	Status       int
	CreateType   int
	ValidType    int
	StartDate    time.Time
	EndDate      time.Time
	Desc         string
	Favicon      string
	ClickNum     int
	TotalPv      int
	TotalUv      int
	TotalUip     int
	TodayPv      int
	TodayUv      int
	TodayUip     int
}
