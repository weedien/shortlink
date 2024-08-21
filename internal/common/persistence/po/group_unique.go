package po

const TableNameGroupUnique = "group_unique"

// GroupUnique mapped from table <group_unique>
type GroupUnique struct {
	ID  int64  `gorm:"column:id;primaryKey;comment:ID" json:"id"`   // ID
	Gid string `gorm:"column:gid;not null;comment:分组标识" json:"gid"` // 分组标识
}

// TableName GroupUnique's table name
func (*GroupUnique) TableName() string {
	return TableNameGroupUnique
}
