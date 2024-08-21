package readrepo

import (
	"context"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"shortlink/internal/common/types"
	"shortlink/internal/recycle_bin/adapters/readrepo/dao"
	"shortlink/internal/recycle_bin/app/query"
)

type RecycleBinQuery struct {
	linkDao dao.LinkDao
}

func NewRecycleBinQuery(db *gorm.DB) RecycleBinQuery {
	return RecycleBinQuery{linkDao: dao.NewLinkDao(db)}
}

// PageDisabledLink 分页统计回收站中的短链接
func (q RecycleBinQuery) PageDisabledLink(
	ctx context.Context,
	param query.PageRecycleBin,
) (res *types.PageResp[query.LinkQueryDTO], err error) {

	r, err := q.linkDao.PageDisabledLink(ctx, param.GidList, param.Current, param.Size)
	if err != nil {
		return
	}

	res = types.ConvertRecords(r, func(dto dao.LinkDTO) query.LinkQueryDTO {
		var queryDTO query.LinkQueryDTO
		if err := copier.Copy(&dto, &queryDTO); err != nil {
			return query.LinkQueryDTO{}
		}
		return queryDTO
	})

	return
}

//func (q RecycleBinQuery) linkStatsToQueryDTOBatch(linkStats []LinkWithStats) (links []query.LinkQueryDTO) {
//
//	for _, linkStat := range linkStats {
//		link := query.LinkQueryDTO{
//			ID:            linkStat.ID,
//			Domain:        linkStat.Domain,
//			ShortUri:      linkStat.ShortUri,
//			FullShortUrl:  linkStat.FullShortUrl,
//			OriginalUrl:   linkStat.OriginalUrl,
//			ClickNum:      int(linkStat.ClickNum),
//			Gid:           linkStat.Gid,
//			EnableStatus:  linkStat.EnableStatus,
//			CreateType:    int(linkStat.CreatedType),
//			ValidDateType: int(linkStat.ValidDateType),
//			ValidDate:     linkStat.ValidDate,
//			Desc:          linkStat.Describe,
//			Favicon:       linkStat.Favicon,
//			TotalPv:       int(linkStat.TotalPv),
//			TotalUv:       int(linkStat.TotalUv),
//			TotalUip:      int(linkStat.TotalUip),
//			TodayPv:       linkStat.TodayPv,
//			TodayUv:       linkStat.TodayUv,
//			TodayUip:      linkStat.TodayUip,
//		}
//		links = append(links, link)
//	}
//	return links
//}
