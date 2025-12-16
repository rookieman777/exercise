package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型（多级关联）
type Comment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Rating    int            `gorm:"default:5;check:rating>=1 AND rating<=5" json:"rating"` // 评分，约束1-5
	Likes     int            `gorm:"default:0" json:"likes"`                                // 点赞数
	PostID    uint           `gorm:"index;not null" json:"post_id"`                         // 外键，指向文章
	UserID    uint           `gorm:"index;not null" json:"user_id"`                         // 外键，指向用户
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`                      // 父评论ID，支持嵌套评论
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Post    Post      `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"`
	User    User      `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Parent  *Comment  `gorm:"foreignKey:ParentID;references:ID" json:"parent,omitempty"`
	Replies []Comment `gorm:"foreignKey:ParentID;references:ID" json:"replies,omitempty"`
}

// TableName 自定义表名
func (Comment) TableName() string {
	return "comments"
}

// BeforeCreate 创建前的钩子
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	// 验证评分范围
	if c.Rating < 1 {
		c.Rating = 1
	} else if c.Rating > 5 {
		c.Rating = 5
	}
	return nil
}

// IsReply 判断是否为回复评论
func (c *Comment) IsReply() bool {
	return c.ParentID != nil
}
