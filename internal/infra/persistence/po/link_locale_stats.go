package po

import (
	"time"
)

const TableNameLinkLocaleStats = "link_locale_stats"

// LinkLocaleStats mapped from table <link_locale_stats>
type LinkLocaleStats struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                                           // 日期
	Cnt          int       `gorm:"column:cnt;comment:访问量" json:"cnt"`                                            // 访问量
	Province     string    `gorm:"column:province;comment:省份名称" json:"province"`                                 // 省份名称
	City         string    `gorm:"column:city;comment:市名称" json:"city"`                                          // 市名称
	Adcode       string    `gorm:"column:adcode;comment:城市编码" json:"adcode"`                                     // 城市编码
	Country      string    `gorm:"column:country;comment:国家标识" json:"country"`                                   // 国家标识
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName LinkLocaleStat's table name
func (*LinkLocaleStats) TableName() string {
	return TableNameLinkLocaleStats
}
