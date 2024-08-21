package po

import (
	"time"
)

const TableNameLinkOsStats = "link_os_stats"

// LinkOsStats mapped from table <link_os_stats>
type LinkOsStats struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                                           // 日期
	Cnt          int       `gorm:"column:cnt;comment:访问量" json:"cnt"`                                            // 访问量
	Os           string    `gorm:"column:os;comment:操作系统" json:"os"`                                             // 操作系统
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkOsStat's table name
func (*LinkOsStats) TableName() string {
	return TableNameLinkOsStats
}
