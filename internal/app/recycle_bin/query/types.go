package query

import "time"

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
	Favicon       string
	TotalPv       int
	TotalUv       int
	TotalUip      int
	TodayPv       int
	TodayUv       int
	TodayUip      int
}
