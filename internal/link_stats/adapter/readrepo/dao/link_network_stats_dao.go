package dao

import (
	"context"
	"gorm.io/gorm"
	"shortlink/internal/user/adapter/po"
)

type LinkNetworkStatsDao struct {
	db *gorm.DB
}

func NewLinkNetworkStatsDao(db *gorm.DB) LinkNetworkStatsDao {
	return LinkNetworkStatsDao{db: db}
}

//// LinkNetworkState 记录访问网络监控数据
//func (m *LinkNetworkStatsDao) LinkNetworkState(linkNetworkStats po.LinkNetworkStats) error {
//	rawSql := `
//INSERT INTO t_link_network_stats (full_short_url, date, cnt, network, create_time, update_time, del_flag)
//VALUES (?, ?, ?, ?, NOW(), NOW(), 0)
//ON DUPLICATE KEY UPDATE cnt = cnt + ?;
//`
//	return m.db.Exec(rawSql, linkNetworkStats.FullShortUrl, linkNetworkStats.Date, linkNetworkStats.Cnt, linkNetworkStats.Network, linkNetworkStats.Cnt).Error
//}

// ListNetworkStatsByLink 根据短链接获取指定日期内访问网络监控数据
func (d *LinkNetworkStatsDao) ListNetworkStatsByLink(ctx context.Context, param LinkQueryParam) ([]po.LinkNetworkStats, error) {
	rawSql := `
SELECT
    tlns.network,
    SUM(tlns.cnt) AS cnt
FROM
    t_link tl INNER JOIN
    t_link_network_stats tlns ON tl.full_short_url = tlns.full_short_url
WHERE
    tlns.full_short_url = ?
    AND tl.gid = ?
    AND tl.del_flag = '0'
    AND tl.enable_status = ?
    AND tlns.date BETWEEN ? and ?
GROUP BY
    tlns.full_short_url, tl.gid, tlns.network;
`
	var result []po.LinkNetworkStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.FullShortUrl, param.Gid, param.EnableStatus, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}

// ListNetworkStatsByGroup 根据分组获取指定日期内访问网络监控数据
func (d *LinkNetworkStatsDao) ListNetworkStatsByGroup(ctx context.Context, param LinkGroupQueryParam) ([]po.LinkNetworkStats, error) {
	rawSql := `
SELECT
    tlns.network,
    SUM(tlns.cnt) AS cnt
FROM
    t_link tl INNER JOIN
    t_link_network_stats tlns ON tl.full_short_url = tlns.full_short_url
WHERE
    tl.gid = ?
    AND tl.del_flag = false
    AND tl.enable_status = '0'
    AND tlns.date BETWEEN ? and ?
GROUP BY
    tl.gid, tlns.network;
`
	var result []po.LinkNetworkStats
	err := d.db.WithContext(ctx).
		Raw(rawSql, param.Gid, param.StartDate, param.EndDate).Scan(&result).Error
	return result, err
}
