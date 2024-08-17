package aggregate

import (
	"shortlink/internal/app/link/query"
	"shortlink/internal/domain/short_link/valobj"
	"shortlink/internal/infra/persistence/po"
	"shortlink/pkg/types"
)

type LinkDTO struct {
	po.Link
	valobj.TodayStatsVo
}

type QueryShortLinkAggregate struct {
	query    query.PageShortLink
	total    int64
	linkDTOs []LinkDTO
}

func NewQueryShortLinkAggregate(query query.PageShortLink) *QueryShortLinkAggregate {
	return &QueryShortLinkAggregate{query: query}
}

func (a QueryShortLinkAggregate) SetLinkDTOs(linkDTOs []LinkDTO) {
	a.linkDTOs = linkDTOs
}

func (a QueryShortLinkAggregate) SetTotal(total int64) {
	a.total = total
}

func (a QueryShortLinkAggregate) GetQuery() query.PageShortLink {
	return a.query
}

// GetPageResp 返回一个PageResp
func (a QueryShortLinkAggregate) GetPageResp() types.PageResp[LinkDTO] {
	return types.PageResp[LinkDTO]{
		Current: a.query.Current,
		Size:    a.query.Size,
		Total:   a.total,
		Records: a.linkDTOs,
	}
}
