package po

import (
	"time"
)

const TableNameLinkDeviceStat = "link_device_stats"

// LinkDeviceStats mapped from table <link_device_stats>
type LinkDeviceStats struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                                           // 日期
	Cnt          int       `gorm:"column:cnt;comment:访问量" json:"cnt"`                                            // 访问量
	Device       string    `gorm:"column:device;comment:设备" json:"device"`                                       // 设备
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkDeviceStat's table name
func (*LinkDeviceStats) TableName() string {
	return TableNameLinkDeviceStat
}
