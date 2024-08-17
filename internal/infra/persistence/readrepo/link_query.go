package readrepo

import (
	"context"
	"github.com/jinzhu/copier"
	"shortlink/common/consts"
	"shortlink/common/types"
	"shortlink/internal/app/link/query"
	"shortlink/internal/infra/persistence/readrepo/internal/dao"
)

type LinkQuery struct {
	linkDao dao.LinkDao
}

func NewLinkQuery(linkDao dao.LinkDao) LinkQuery {
	return LinkQuery{linkDao: linkDao}
}

func (q LinkQuery) PageLink(ctx context.Context, param query.PageLink) (res *types.PageResp[query.Link], err error) {
	linkPage, err := q.linkDao.PageLink(param.Gid, consts.StatusEnable, param.OrderTag, param.Current, param.Size)

	if err != nil {
		return
	}

	res = types.ConvertRecords(linkPage, func(linkDTO dao.LinkDTO) query.Link {
		var queryLink query.Link
		err := copier.Copy(&linkDTO, &queryLink)
		if err != nil {
			return query.Link{}
		}
		return queryLink
	})
	return
}

func (q LinkQuery) ListGroupLinkCount(ctx context.Context, gidList []string) (res []query.GroupLinkCount, err error) {
	linkGidCountDTO, err := q.linkDao.ListGroupLinkCount(gidList)
	if err != nil {
		return
	}

	res = make([]query.GroupLinkCount, 0)
	for _, dto := range linkGidCountDTO {
		res = append(res, query.GroupLinkCount{
			Gid:   dto.Gid,
			Count: dto.Count,
		})
	}
	return
}
