package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/common/persistence/po"
)

type LinkOsStatsDao struct {
	db *gorm.DB
}

func NewLinkOsStatsDao(db *gorm.DB) LinkOsStatsDao {
	return LinkOsStatsDao{db: db}
}

//// LinkOsState 记录操作系统访问监控数据
//func (m *LinkOsStatsDao) LinkOsState(linkOsStats po.LinkOsStats) error {
//	rawSql := `
//INSERT INTO t_link_os_stats (full_short_url, date, cnt, os, create_time, update_time, del_flag)
//VALUES (?, ?, ?, ?, NOW(), NOW(), 0)
//ON DUPLICATE KEY UPDATE cnt = cnt + ?;
//`
//	return m.db.Exec(rawSql, linkOsStats.FullShortUrl, linkOsStats.Date, linkOsStats.Cnt, linkOsStats.Os, linkOsStats.Cnt).Error
//}

// ListOsStatsByLink 根据短链接获取指定日期内操作系统监控数据
func (d *LinkOsStatsDao) ListOsStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkOsStats, error) {
	rawSql := `
SELECT
    tlos.os,
    SUM(tlos.cnt) AS count
FROM
    t_link tl INNER JOIN
    t_link_os_stats tlos ON tl.full_short_url = tlos.full_short_url
WHERE
    tlos.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlos.date BETWEEN ? and ?
GROUP BY
    tlos.full_short_url, tl.gid, tlos.os;
`
	var result []po.LinkOsStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListOsStatsByGroup 根据分组获取指定日期内操作系统监控数据
func (d *LinkOsStatsDao) ListOsStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkOsStats, error) {
	rawSql := `
SELECT
    tlos.os,
    SUM(tlos.cnt) AS count
FROM
    t_link tl INNER JOIN
    t_link_os_stats tlos ON tl.full_short_url = tlos.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlos.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlos.os;
`
	var result []po.LinkOsStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}
