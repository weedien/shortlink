package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/common/persistence/po"
)

type LinkBrowserStatsDao struct {
	db *gorm.DB
}

func NewLinkBrowserStatsDao(db *gorm.DB) LinkBrowserStatsDao {
	return LinkBrowserStatsDao{db: db}
}

//// LinkBrowserState 记录浏览器访问监控数据
//func (m *LinkBrowserStatsDao) LinkBrowserState(linkBrowserStats po.LinkBrowserStats) error {
//	rawSql := `
//INSERT INTO t_link_browser_stats (full_short_url, date, cnt, browser, create_time, update_time, del_flag)
//VALUES (?, ?, ?, ?, NOW(), NOW(), 0)
//ON DUPLICATE KEY UPDATE cnt = cnt + ?;
//`
//	return m.db.Exec(rawSql, linkBrowserStats.FullShortUrl, linkBrowserStats.Date, linkBrowserStats.Cnt, linkBrowserStats.Browser, linkBrowserStats.Cnt).Error
//}

// ListBrowserStatsByLink 根据短链接获取指定日期内浏览器监控数据
func (d *LinkBrowserStatsDao) ListBrowserStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkBrowserStats, error) {
	rawSql := `
SELECT
    tlbs.browser,
    SUM(tlbs.cnt) AS cnt
FROM
    t_link tl INNER JOIN
    t_link_browser_stats tlbs ON tl.full_short_url = tlbs.full_short_url
WHERE
    tlbs.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlbs.date BETWEEN ? and ?
GROUP BY
    tlbs.full_short_url, tl.gid, tlbs.browser;
`
	var result []po.LinkBrowserStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListBrowserStatsByGroup 根据分组获取指定日期内浏览器监控数据
func (d *LinkBrowserStatsDao) ListBrowserStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkBrowserStats, error) {
	rawSql := `
SELECT
    tlbs.browser,
    SUM(tlbs.cnt) AS count
FROM
    t_link tl INNER JOIN
    t_link_browser_stats tlbs ON tl.full_short_url = tlbs.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = '0'
    AND tlbs.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlbs.browser;
`
	var result []po.LinkBrowserStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}
