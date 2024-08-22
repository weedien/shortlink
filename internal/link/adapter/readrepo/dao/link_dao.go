package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/common/types"
	po2 "shortlink/internal/link/adapter/po"
	"shortlink/internal/link/domain/entity"
)

type LinkDao struct {
	db *gorm.DB
}

func NewLinkDao(db *gorm.DB) LinkDao {
	return LinkDao{db: db}
}

//// IncrementStats 短链接访问统计自增
//func (m *LinkDao) IncrementStats(gid string, param.FullShortUrl string, totalPv int, totalUv int, totalUip int) error {
//	rawSql := `
//UPDATE t_link
//SET total_pv = total_pv + ?, total_uv = total_uv + ?, total_uip = total_uip + ?
//WHERE gid = ? AND full_short_url = ?;
//`
//	return m.db.Exec(rawSql, totalPv, totalUv, totalUip, gid, param.FullShortUrl).Error
//}

type LinkDTO struct {
	po2.Link
	TodayPv  int `json:"todayPv"`
	TodayUv  int `json:"todayUv"`
	TodayUip int `json:"todayUip"`
}

type LinkGidCountDTO struct {
	Gid   string `json:"gid"`
	Count int    `json:"count"`
}

func (d *LinkDao) GetLink(ctx context.Context, id entity.LinkID) (po2.Link, error) {
	var link po2.Link
	err := d.db.WithContext(ctx).
		Where("full_short_url = ? AND gid = ? AND enable_status = 0 AND del_flag = false", id.FullShortUrl, id.Gid).
		First(&link).Error
	return link, err
}

func (d *LinkDao) GetLinkGoto(ctx context.Context, fullShortUrl string) (po2.LinkGoto, error) {
	var linkGoto po2.LinkGoto
	err := d.db.WithContext(ctx).
		Where("full_short_url = ?", fullShortUrl).
		First(&linkGoto).Error
	return linkGoto, err

}

// ListGroupLinkCount 统计分组短链接数量
func (d *LinkDao) ListGroupLinkCount(ctx context.Context, gidList []string) (res []LinkGidCountDTO, err error) {
	err = d.db.WithContext(ctx).
		Table("link").
		Select("gid, COUNT(*) AS count").
		Where("gid IN ?", gidList).
		Where("enable_status = 0 AND del_flag = false").
		Group("gid").
		Find(&res).Error
	return
}

// PageLink 分页统计短链接
func (d *LinkDao) PageLink(
	ctx context.Context,
	gid string,
	enableStatus int,
	orderTag string,
	page int,
	size int,
) (res *types.PageResp[LinkDTO], err error) {
	rawSql := `
SELECT t.*, COALESCE(s.today_pv, 0) AS todayPv, COALESCE(s.today_uv, 0) AS todayUv, COALESCE(s.today_uip, 0) AS todayUip
FROM t_link t
LEFT JOIN t_link_stats_today s ON t.full_short_url = s.full_short_url AND s.date = current_date
WHERE t.gid = ? AND t.enable_status = ? AND t.del_flag = 0
ORDER BY 
    CASE 
        WHEN ? = 'todayPv' THEN todayPv
        WHEN ? = 'todayUv' THEN todayUv
        WHEN ? = 'todayUip' THEN todayUip
        WHEN ? = 'totalPv' THEN t.total_pv
        WHEN ? = 'totalUv' THEN t.total_uv
        WHEN ? = 'totalUip' THEN t.total_uip
        ELSE t.create_time
    END DESC
LIMIT ? OFFSET ?;
`

	var records []LinkDTO
	err = d.db.WithContext(ctx).
		Raw(rawSql, gid, enableStatus, orderTag, orderTag, orderTag, orderTag, orderTag, orderTag, size, (page-1)*size).Scan(&records).Error
	if err != nil {
		return
	}

	var total int64
	err = d.db.WithContext(ctx).
		Model(&po2.Link{}).Where("gid = ? AND enable_status = ? AND del_flag = 0", gid, enableStatus).Count(&total).Error
	if err != nil {
		return
	}

	res = &types.PageResp[LinkDTO]{
		Current: page,
		Size:    size,
		Total:   total,
		Records: records,
	}
	return
}

// PageDisabledLink 分页统计回收站短链接
func (d *LinkDao) PageDisabledLink(
	ctx context.Context,
	gidList []string,
	page int,
	size int,
) (*types.PageResp[LinkDTO], error) {
	rawSql := `
SELECT t.*, COALESCE(s.today_pv, 0) AS todayPv, COALESCE(s.today_uv, 0) AS todayUv, COALESCE(s.today_uip, 0) AS todayUip
FROM t_link t
LEFT JOIN t_link_stats_today s ON t.full_short_url = s.full_short_url AND s.date = current_date
WHERE t.gid IN (?) AND t.enable_status = 1 AND t.del_flag = 0
ORDER BY t.update_time
LIMIT ? OFFSET ?;
`
	var records []LinkDTO
	err := d.db.WithContext(ctx).
		Raw(rawSql, gidList, size, (page-1)*size).Scan(&records).Error

	if err != nil {
		return nil, err
	}

	var total int64
	err = d.db.WithContext(ctx).
		Model(&po2.Link{}).Where("gid IN (?) AND enable_status = 1 AND del_flag = 0", gidList).Count(&total).Error
	if err != nil {
		return nil, err
	}

	res := &types.PageResp[LinkDTO]{
		Current: page,
		Size:    size,
		Total:   total,
		Records: records,
	}

	return res, err
}
