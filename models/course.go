package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Course 课程模型（多对多关联）
type Course struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"` // 课程名称
	Title       string         `gorm:"type:varchar(255);unique;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Code        string         `gorm:"type:varchar(20);unique;not null" json:"code"` // 课程代码
	Category    string         `gorm:"type:varchar(50)" json:"category"`             // 课程分类
	Price       float64        `gorm:"type:decimal(10,2);default:0.00;check:price>=0" json:"price"`
	Duration    int            `gorm:"default:40" json:"duration"` // 课时数
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Tags        datatypes.JSON `gorm:"type:json" json:"tags"`     // JSON存储标签
	Schedule    datatypes.JSON `gorm:"type:json" json:"schedule"` // JSON存储课程安排
	StartDate   time.Time      `gorm:"index" json:"start_date"`
	EndDate     time.Time      `gorm:"index" json:"end_date"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	// 多对多关联
	Users    []User `gorm:"many2many:user_courses;" json:"users,omitempty"`
	Teachers []User `gorm:"many2many:course_teachers;" json:"teachers,omitempty"`
}

// TableName 自定义表名
func (Course) TableName() string {
	return "courses"
}

// BeforeSave 保存前的钩子
func (c *Course) BeforeSave(tx *gorm.DB) error {
	// 验证日期逻辑
	if !c.EndDate.IsZero() && c.StartDate.After(c.EndDate) {
		return fmt.Errorf("开始日期不能晚于结束日期")
	}
	return nil
}

// IsOngoing 检查课程是否正在进行中
func (c *Course) IsOngoing() bool {
	now := time.Now()
	return c.IsActive && now.After(c.StartDate) && (c.EndDate.IsZero() || now.Before(c.EndDate))
}

// UserCourse 用户选课中间表模型
type UserCourse struct {
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	CourseID   uint      `gorm:"primaryKey" json:"course_id"`
	EnrolledAt time.Time `gorm:"autoCreateTime" json:"enrolled_at"`
	Grade      *float64  `gorm:"type:decimal(5,2)" json:"grade,omitempty"` // 成绩，可为空
	Status     string    `gorm:"type:enum('enrolled','completed','dropped');default:'enrolled'" json:"status"`

	// 关联关系
	User   User   `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	Course Course `gorm:"foreignKey:CourseID;references:ID" json:"course,omitempty"`
}

// TableName 自定义中间表名
func (UserCourse) TableName() string {
	return "user_courses"
}
