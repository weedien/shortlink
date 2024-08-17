package aggregate

import (
	"shortlink/internal/domain/link_status/valobj"
)

type LinkStatsAggregate struct {
	// 入参
	params valobj.ShortLinkSimpleVO

	// 出参
	result *valobj.ShortLinkStats

	// 异常
	err error
}

func NewLinkStatsAggregate(params valobj.ShortLinkSimpleVO) *LinkStatsAggregate {
	return &LinkStatsAggregate{params: params}
}

func (a LinkStatsAggregate) Result() (*valobj.ShortLinkStats, error) {
	return a.result, a.err
}
