package po

import "gorm.io/gorm"

const TableNameLinkGoto = "link_goto"

// LinkGoto mapped from table <link_goto>
type LinkGoto struct {
	ID           int            `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"`       // ID
	Gid          string         `gorm:"column:gid;not null;comment:分组标识" json:"gid"`                        // 分组标识
	FullShortUrl string         `gorm:"column:full_short_url;not null;comment:完整短链接" json:"full_short_url"` // 完整短链接
	DeleteTime   gorm.DeletedAt `gorm:"column:delete_time;comment:删除时间" json:"delete_time"`                 // 删除时间
}

// TableName LinkGoto's table name
func (*LinkGoto) TableName() string {
	return TableNameLinkGoto
}
