package po

import (
	"gorm.io/gorm"
	"time"
)

const TableNameLink = "link"

// Link mapped from table <link>
type Link struct {
	ID            int            `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	Domain        string         `gorm:"column:domain;comment:域名" json:"domain"`                                       // 域名
	ShortUri      string         `gorm:"column:short_uri;comment:短链接" json:"short_uri"`                                // 短链接
	FullShortUrl  string         `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`                    // 完整短链接
	OriginalUrl   string         `gorm:"column:origin_url;comment:原始链接" json:"origin_url"`                             // 原始链接
	ClickNum      int            `gorm:"column:click_num;comment:点击量" json:"click_num"`                                // 点击量
	Gid           string         `gorm:"column:gid;comment:分组标识" json:"gid"`                                           // 分组标识
	Favicon       string         `gorm:"column:favicon;comment:网站图标" json:"favicon"`                                   // 网站图标
	EnableStatus  int            `gorm:"column:enable_status;comment:启用标识 0：启用 1：未启用" json:"enable_status"`            // 启用标识 0：启用 1：未启用
	CreatedType   int            `gorm:"column:created_type;comment:创建类型 0：接口创建 1：控制台创建" json:"created_type"`          // 创建类型 0：接口创建 1：控制台创建
	ValidDateType int            `gorm:"column:valid_date_type;comment:有效期类型 0：永久有效 1：自定义" json:"valid_date_type"`     // 有效期类型 0：永久有效 1：自定义
	ValidDate     time.Time      `gorm:"column:valid_date;comment:有效期" json:"valid_date"`                              // 有效期
	Describe      string         `gorm:"column:describe;comment:描述" json:"describe"`                                   // 描述
	TotalPv       int            `gorm:"column:total_pv;comment:历史PV" json:"total_pv"`                                 // 历史PV
	TotalUv       int            `gorm:"column:total_uv;comment:历史UV" json:"total_uv"`                                 // 历史UV
	TotalUip      int            `gorm:"column:total_uip;comment:历史UIP" json:"total_uip"`                              // 历史UIP
	CreateTime    time.Time      `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime    time.Time      `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelTime       gorm.DeletedAt `gorm:"column:del_time;comment:删除时间戳" json:"del_time"`                                // 删除时间戳
}

// TableName Link's table name
func (*Link) TableName() string {
	return TableNameLink
}

func (l *Link) BeforeSave(tx *gorm.DB) (err error) {
	if l.ValidDateType == 0 {
		tx.Statement.SetColumn("valid_date", nil)
	}
	return
}

func (l *Link) BeforeUpdate(tx *gorm.DB) (err error) {
	if l.ValidDateType == 0 {
		tx.Statement.SetColumn("valid_date", nil)
	}
	return
}