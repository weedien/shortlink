package dao

import "time"

type LinkQueryParam struct {
	FullShortUrl string
	Gid          string
	EnableStatus int
	StartDate    time.Time
	EndDate      time.Time
}

type LinkGroupQueryParam struct {
	Gid          string
	EnableStatus int
	StartDate    time.Time
	EndDate      time.Time
}
