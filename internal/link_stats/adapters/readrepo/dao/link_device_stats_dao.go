package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/common/persistence/po"
)

type LinkDeviceStatsDao struct {
	db *gorm.DB
}

func NewLinkDeviceStatsDao(db *gorm.DB) LinkDeviceStatsDao {
	return LinkDeviceStatsDao{db: db}
}

//// LinkDeviceState 记录访问设备监控数据
//func (m *LinkDeviceStatsDao) LinkDeviceState(linkDeviceStats po.LinkDeviceStats) error {
//	rawSql := `
//INSERT INTO t_link_device_stats (full_short_url, date, cnt, device, create_time, update_time, del_flag)
//VALUES (?, ?, ?, ?, NOW(), NOW(), 0)
//ON DUPLICATE KEY UPDATE cnt = cnt + ?;
//`
//	return m.db.Exec(rawSql, linkDeviceStats.FullShortUrl, linkDeviceStats.Date, linkDeviceStats.Cnt, linkDeviceStats.Device, linkDeviceStats.Cnt).Error
//}

// ListDeviceStatsByLink 根据短链接获取指定日期内访问设备监控数据
func (d *LinkDeviceStatsDao) ListDeviceStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkDeviceStats, error) {
	rawSql := `
SELECT
    tlds.device,
    SUM(tlds.cnt) AS cnt
FROM
    t_link tl INNER JOIN
    t_link_device_stats tlds ON tl.full_short_url = tlds.full_short_url
WHERE
    tlds.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlds.date BETWEEN ? and ?
GROUP BY
    tlds.full_short_url, tl.gid, tlds.device;
`
	var result []po.LinkDeviceStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListDeviceStatsByGroup 根据分组获取指定日期内访问设备监控数据
func (d *LinkDeviceStatsDao) ListDeviceStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkDeviceStats, error) {
	rawSql := `
SELECT
    tlds.device,
    SUM(tlds.cnt) AS cnt
FROM
    t_link tl INNER JOIN
    t_link_device_stats tlds ON tl.full_short_url = tlds.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlds.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlds.device;
`
	var result []po.LinkDeviceStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}
