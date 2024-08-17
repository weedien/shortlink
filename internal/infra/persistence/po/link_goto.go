package po

const TableNameLinkGoto = "link_goto"

// LinkGoto mapped from table <link_goto>
type LinkGoto struct {
	ID           int64  `gorm:"column:id;primaryKey;comment:ID" json:"id"`                 // ID
	Gid          string `gorm:"column:gid;comment:分组标识" json:"gid"`                        // 分组标识
	FullShortURL string `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"` // 完整短链接
}

// TableName LinkGoto's table name
func (*LinkGoto) TableName() string {
	return TableNameLinkGoto
}
