package aggregate

import "shortlink/internal/link_stats/domain/valobj"

type LinkStatAggregate struct {
	// 入参
	params valobj.ShortLinkSimpleVO

	// 出参
	result *valobj.ShortLinkStat

	// 异常
	err error
}

func NewLinkStatAggregate(params valobj.ShortLinkSimpleVO) *LinkStatAggregate {
	return &LinkStatAggregate{params: params}
}

func (a LinkStatAggregate) Result() (*valobj.ShortLinkStat, error) {
	return a.result, a.err
}
