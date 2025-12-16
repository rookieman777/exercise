package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Post 文章模型（一对多关联）
type Post struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null;index" json:"title"`                           // 标题，创建索引
	Content     string         `gorm:"type:text;not null" json:"content"`                                       // 内容
	Slug        string         `gorm:"type:varchar(255);unique;not null" json:"slug"`                           // URL友好标识
	Status      string         `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"` // 枚举类型
	Views       int            `gorm:"default:0" json:"views"`                                                  // 浏览量
	Tags        datatypes.JSON `gorm:"type:json" json:"tags"`                                                   // JSON字段存储标签
	PublishedAt *time.Time     `gorm:"index" json:"published_at,omitempty"`                                     // 发布时间，索引（使用指针允许NULL）
	AuthorID    uint           `gorm:"index" json:"author_id"`                                                  // 外键，指向作者
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Author   User      `gorm:"foreignKey:AuthorID;references:ID" json:"author,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments,omitempty"`
}

// TableName 自定义表名
func (Post) TableName() string {
	return "posts"
}

// BeforeSave 保存前的钩子
func (p *Post) BeforeSave(tx *gorm.DB) error {
	// 自动设置发布时间
	if p.Status == "published" && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

// IncrementViews 增加浏览量
func (p *Post) IncrementViews() {
	p.Views++
}

// IsPublished 检查文章是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == "published" && p.PublishedAt != nil && !p.PublishedAt.After(time.Now())
}
