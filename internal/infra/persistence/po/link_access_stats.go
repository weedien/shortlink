package po

import (
	"time"
)

const TableNameLinkAccessStat = "link_access_stats"

// LinkAccessStats mapped from table <link_access_stats>
type LinkAccessStats struct {
	ID           int       `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                                           // 日期
	Pv           int       `gorm:"column:pv;comment:访问量" json:"pv"`                                              // 访问量
	Uv           int       `gorm:"column:uv;comment:独立访客数" json:"uv"`                                            // 独立访客数
	Uip          int       `gorm:"column:uip;comment:独立IP数" json:"uip"`                                          // 独立IP数
	Hour         int       `gorm:"column:hour;comment:小时" json:"hour"`                                           // 小时
	Weekday      int       `gorm:"column:weekday;comment:星期" json:"weekday"`                                     // 星期
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkAccessStat's table name
func (*LinkAccessStats) TableName() string {
	return TableNameLinkAccessStat
}
