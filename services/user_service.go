package services

import (
	//"encoding/json"
	"errors"
	"fmt"

	//"log"
	//"time"

	"exercise/models"
	"exercise/repositories"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound   = errors.New("用户不存在")
	ErrInvalidEmail   = errors.New("邮箱格式错误")
	ErrInvalidAge     = errors.New("年龄必须大于0且小于150")
	ErrWeakPassword   = errors.New("密码强度不足")
	ErrDuplicateEmail = errors.New("邮箱已存在")
)

// UserService 用户服务接口
type UserService interface {
	Register(user *models.User) error
	// Login(email, password string) (*models.User, error)
	// GetUserByID(id uint) (*models.User, error)
	// UpdateProfile(id uint, profile *models.Profile) error
	// DeactivateAccount(id uint) error
	// SearchUsers(keyword string, page, pageSize int) ([]models.User, int64, error)
	// GetUserStats() (*UserStats, error)
	// ExportUsers() ([]byte, error)
	// ImportUsers(data []byte) (int, error)
}

// UserStats 用户统计信息
type UserStats struct {
	TotalUsers     int64    `json:"total_users"`
	ActiveUsers    int64    `json:"active_users"`
	InactiveUsers  int64    `json:"inactive_users"`
	TodayRegisters int64    `json:"today_registers"`
	AvgAge         float64  `json:"avg_age"`
	TopDomains     []string `json:"top_domains"`
}

// userServiceImpl 用户服务实现
type userServiceImpl struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userServiceImpl{
		userRepo: repositories.NewUserRepository(),
	}
}

// Register 用户注册
func (s *userServiceImpl) Register(user *models.User) error {
	// 验证逻辑
	if user.Age <= 0 || user.Age > 150 {
		return ErrInvalidAge
	}

	if len(user.Password) < 8 {
		return ErrWeakPassword
	}

	// 检查邮箱是否已存在
	existing, err := s.userRepo.FindByEmail(user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查邮箱失败: %v", err)
	}
	if existing != nil {
		return ErrDuplicateEmail
	}

	// 创建用户（自动创建关联profile）
	return s.userRepo.Create(user)
}
