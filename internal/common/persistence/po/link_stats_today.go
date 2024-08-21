package po

import (
	"time"
)

const TableNameLinkStatsToday = "link_stats_today"

// LinkStatsToday mapped from table <link_stats_today>
type LinkStatsToday struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:短链接" json:"full_short_url"`                      // 短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                                           // 日期
	TodayPv      int       `gorm:"column:today_pv;comment:今日PV" json:"today_pv"`                                 // 今日PV
	TodayUv      int       `gorm:"column:today_uv;comment:今日UV" json:"today_uv"`                                 // 今日UV
	TodayUip     int       `gorm:"column:today_uip;comment:今日IP数" json:"today_uip"`                              // 今日IP数
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkStatsToday's table name
func (*LinkStatsToday) TableName() string {
	return TableNameLinkStatsToday
}
