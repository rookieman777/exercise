package models

import (
	"time"
)

// Profile 用户资料模型（一对一关联）
type Profile struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"unique;not null" json:"user_id"` // 外键，与User一对一
	FirstName string    `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50);not null" json:"last_name"`
	Bio       string    `gorm:"type:text" json:"bio"`
	AvatarURL string    `gorm:"type:varchar(255)" json:"avatar_url"`
	Location  string    `gorm:"type:varchar(100)" json:"location"`
	Website   string    `gorm:"type:varchar(255)" json:"website"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// TableName 自定义表名
func (Profile) TableName() string {
	return "profiles"
}

// FullName 计算全名
func (p *Profile) FullName() string {
	return p.FirstName + " " + p.LastName
}
