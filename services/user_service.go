package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"log"
	"time"

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
	Login(email, password string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateProfile(id uint, profile *models.Profile) error
	DeactivateAccount(id uint) error
	SearchUsers(keyword string, page, pageSize int) ([]models.User, int64, error)
	GetUserStats() (*UserStats, error)
	ExportUsers() ([]byte, error)
	ImportUsers(data []byte) (int, error)
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

// Login 用户登录
func (s *userServiceImpl) Login(email, password string) (*models.User, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查找用户失败: %v", err)
	}

	// 验证密码（实际应用中应该使用bcrypt等加密方式）
	if user.Password != password {
		return nil, errors.New("密码错误")
	}

	// TODO: 更新最后登录时间（需要在User模型中添加LastLoginAt字段）
	// user.LastLoginAt = time.Now()
	// if err := s.userRepo.Update(user); err != nil {
	// 	log.Printf("更新登录时间失败: %v", err)
	// }

	return user, nil
}

// GetUserByID 获取用户详情（包含关联数据）
func (s *userServiceImpl) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("获取用户失败: %v", err)
	}

	// TODO: 获取用户的文章和评论统计（需要在User模型中添加Extra字段）
	// postCount, commentCount, err := s.getUserActivityStats(id)
	// if err != nil {
	// 	log.Printf("获取用户活动统计失败: %v", err)
	// }
	//
	// user.Extra = map[string]interface{}{
	// 	"post_count":    postCount,
	// 	"comment_count": commentCount,
	// }

	return user, nil
}

// UpdateProfile 更新用户资料
func (s *userServiceImpl) UpdateProfile(id uint, profile *models.Profile) error {
	// 获取用户
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	// 更新profile
	if user.Profile == nil {
		user.Profile = profile
	} else {
		user.Profile = mergeProfiles(user.Profile, profile)
	}

	// 保存更新
	return s.userRepo.Update(user)
}

// DeactivateAccount 停用账户（软删除）
func (s *userServiceImpl) DeactivateAccount(id uint) error {
	// 停用用户
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	user.IsActive = false
	// TODO: 添加DeactivatedAt字段到User模型
	// user.DeactivatedAt = time.Now()

	return s.userRepo.Update(user)
}

// SearchUsers 搜索用户（分页）
func (s *userServiceImpl) SearchUsers(keyword string, page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.Search(keyword, page, pageSize)
}

// GetUserStats 获取用户统计信息
func (s *userServiceImpl) GetUserStats() (*UserStats, error) {
	var stats UserStats

	// 获取所有用户进行统计
	users, _, err := s.userRepo.FindAll(1, 1000000)
	if err != nil {
		return nil, fmt.Errorf("获取用户数据失败: %v", err)
	}

	stats.TotalUsers = int64(len(users))
	var totalAge int
	for _, user := range users {
		if user.IsActive {
			stats.ActiveUsers++
		} else {
			stats.InactiveUsers++
		}
		totalAge += user.Age

		// 统计今日注册
		if user.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			stats.TodayRegisters++
		}
	}

	if stats.TotalUsers > 0 {
		stats.AvgAge = float64(totalAge) / float64(stats.TotalUsers)
	}

	return &stats, nil
}

// ExportUsers 导出用户数据（JSON格式）
func (s *userServiceImpl) ExportUsers() ([]byte, error) {
	users, _, err := s.userRepo.FindAll(1, 1000000) // 获取所有用户
	if err != nil {
		return nil, fmt.Errorf("获取用户数据失败: %v", err)
	}

	// 转换为JSON
	return json.Marshal(users)
}

// ImportUsers 导入用户数据（批量创建）
func (s *userServiceImpl) ImportUsers(data []byte) (int, error) {
	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return 0, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 批量创建
	count := 0
	for _, user := range users {
		if err := s.userRepo.Create(&user); err != nil {
			// 跳过重复邮箱的用户
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				continue
			}
			log.Printf("导入用户失败: %v", err)
			continue
		}
		count++
	}

	return count, nil
}

// 辅助函数：合并用户资料
func mergeProfiles(old, new *models.Profile) *models.Profile {
	result := *old
	if new.Bio != "" {
		result.Bio = new.Bio
	}
	if new.Location != "" {
		result.Location = new.Location
	}
	if new.AvatarURL != "" {
		result.AvatarURL = new.AvatarURL
	}
	if new.Website != "" {
		result.Website = new.Website
	}
	return &result
}
