package po

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID           int64     `gorm:"column:id;primaryKey;comment:ID" json:"id"`                                    // ID
	Username     string    `gorm:"column:username;not null;comment:用户名" json:"username"`                         // 用户名
	Password     string    `gorm:"column:password;comment:密码" json:"password"`                                   // 密码
	RealName     string    `gorm:"column:real_name;comment:真实姓名" json:"real_name"`                               // 真实姓名
	Phone        string    `gorm:"column:phone;comment:手机号" json:"phone"`                                        // 手机号
	Mail         string    `gorm:"column:mail;comment:邮箱" json:"mail"`                                           // 邮箱
	DeletionTime time.Time `gorm:"column:deletion_time;comment:注销时间戳" json:"deletion_time"`                      // 注销时间戳
	CreateTime   time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;comment:创建时间" json:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;comment:修改时间" json:"update_time"` // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`                     // 删除标识 0：未删除 1：已删除
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
