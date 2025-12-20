package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`               // 主键，自增
	Username  string         `gorm:"type:varchar(50);unique;not null" json:"username"` // 用户名，唯一，非空
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`   // 邮箱，唯一，非空
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`              // 密码，不返回JSON
	Age       int            `gorm:"default:18;check:age>=0" json:"age"`               // 年龄，默认18，约束>=0
	IsActive  bool           `gorm:"default:true" json:"is_active"`                    // 是否活跃，默认true
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间，自动设置
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                 // 更新时间，自动更新
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                // 软删除时间戳

	// 关联关系
	Profile *Profile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"profile,omitempty"`
	Posts   []Post   `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"posts,omitempty"`
	Courses []Course `gorm:"many2many:user_courses;" json:"courses,omitempty"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 示例：创建前自动设置默认值或验证
	if u.Age == 0 {
		u.Age = 18
	}
	return nil
}

// AfterCreate 创建后的钩子
func (u *User) AfterCreate(tx *gorm.DB) error {
	// 示例：创建后的钩子（已禁用自动创建Profile，避免重复插入）
	// 如需自动创建Profile，请在业务代码中手动处理
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 示例：更新前记录变更日志
	// log.Printf("更新用户: %s (ID: %d)", u.Username, u.ID)
	return nil
}
