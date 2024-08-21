package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/common/persistence/po"
)

type LinkAccessStatsDao struct {
	db *gorm.DB
}

func NewLinkAccessStatsDao(db *gorm.DB) LinkAccessStatsDao {
	return LinkAccessStatsDao{db: db}
}

//// LinkStats 记录基础访问监控数据
//func (m *LinkAccessStatsDao) LinkStats(linkAccessStats po.LinkAccessStats) error {
//	rawSql := `
//INSERT INTO t_link_access_stats (full_short_url, date, pv, uv, uip, hour, weekday, create_time, update_time, del_flag)
//VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), 0)
//ON DUPLICATE KEY UPDATE pv = pv + ?, uv = uv + ?, uip = uip + ?;
//`
//	return m.db.Exec(rawSql, linkAccessStats.FullShortUrl, linkAccessStats.Date, linkAccessStats.Pv, linkAccessStats.Uv, linkAccessStats.Uip, linkAccessStats.Hour, linkAccessStats.Weekday, linkAccessStats.Pv, linkAccessStats.Uv, linkAccessStats.Uip).Error
//}

// ListStatsByLink 根据短链接获取指定日期内基础监控数据
func (d *LinkAccessStatsDao) ListStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.date,
    SUM(tlas.pv) AS pv,
    SUM(tlas.uv) AS uv,
    SUM(tlas.uip) AS uip
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tlas.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tlas.full_short_url, tl.gid, tlas.date;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListStatsByGroup 根据分组获取指定日期内基础监控数据
func (d *LinkAccessStatsDao) ListStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.date,
    SUM(tlas.pv) AS pv,
    SUM(tlas.uv) AS uv,
    SUM(tlas.uip) AS uip
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = '0'
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlas.date;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListHourStatsByLink 根据短链接获取指定日期内小时基础监控数据
func (d *LinkAccessStatsDao) ListHourStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.hour,
    SUM(tlas.pv) AS pv
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tlas.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tlas.full_short_url, tl.gid, tlas.hour;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListHourStatsByGroup 根据分组获取指定日期内小时基础监控数据
func (d *LinkAccessStatsDao) ListHourStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.hour,
    SUM(tlas.pv) AS pv
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = '0'
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlas.hour;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListWeekdayStatsByLink 根据短链接获取指定日期内小时基础监控数据
func (d *LinkAccessStatsDao) ListWeekdayStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.weekday,
    SUM(tlas.pv) AS pv
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tlas.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tlas.full_short_url, tl.gid, tlas.weekday;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListWeekdayStatsByGroup 根据分组获取指定日期内小时基础监控数据
func (d *LinkAccessStatsDao) ListWeekdayStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkAccessStats, error) {
	rawSql := `
SELECT
    tlas.weekday,
    SUM(tlas.pv) AS pv
FROM
    t_link tl INNER JOIN
    t_link_access_stats tlas ON tl.full_short_url = tlas.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = '0'
    AND tlas.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlas.weekday;
`
	var result []po.LinkAccessStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}
