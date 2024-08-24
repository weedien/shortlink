package po

import (
	"gorm.io/gorm"
	"time"
)

const TableNameLinkOsStat = "link_os_stats"

// LinkOsStat mapped from table <link_os_stats>
type LinkOsStat struct {
	ID           int            `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"`                 // ID
	FullShortUrl string         `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	Date         time.Time      `gorm:"column:date;not null;default:CURRENT_DATE;comment:日期" json:"date"`             // 日期
	Cnt          int            `gorm:"column:cnt;comment:访问量" json:"cnt"`                                            // 访问量
	Os           string         `gorm:"column:os;comment:操作系统" json:"os"`                                             // 操作系统
	CreateTime   time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                           // 删除时间
}

// TableName LinkOsStat's table name
func (*LinkOsStat) TableName() string {
	return TableNameLinkOsStat
}
