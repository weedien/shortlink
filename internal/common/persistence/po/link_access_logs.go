package po

import (
	"time"
)

const TableNameLinkAccessLog = "link_access_logs"

// LinkAccessLog mapped from table <link_access_logs>
type LinkAccessLog struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	User         string    `gorm:"column:user;comment:用户信息" json:"user"`                                         // 用户信息
	IP           string    `gorm:"column:ip;comment:IP" json:"ip"`                                               // IP
	Browser      string    `gorm:"column:browser;comment:浏览器" json:"browser"`                                    // 浏览器
	Os           string    `gorm:"column:os;comment:操作系统" json:"os"`                                             // 操作系统
	Network      string    `gorm:"column:network;comment:访问网络" json:"network"`                                   // 访问网络
	Device       string    `gorm:"column:device;comment:访问设备" json:"device"`                                     // 访问设备
	Locale       string    `gorm:"column:locale;comment:地区" json:"locale"`                                       // 地区
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkAccessLog's table name
func (*LinkAccessLog) TableName() string {
	return TableNameLinkAccessLog
}
