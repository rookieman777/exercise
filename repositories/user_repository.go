package repositories

import (
	//"fmt"

	"exercise/database"
	"exercise/models"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	// FindByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	// Delete(id uint) error
	// FindAll(page, pageSize int) ([]models.User, int64, error)
	Search(keyword string, page, pageSize int) ([]models.User, int64, error)
	// Count() (int64, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建新的用户仓储实例
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查找用户（包含关联数据）
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Profile").Preload("Posts").First(&user, id).Error//Preload避免N+1查询
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}


// Search 搜索用户
func (r *userRepository) Search(keyword string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	offset := (page - 1) * pageSize
	searchPattern := "%" + keyword + "%"

	// 构建查询条件
	query := r.db.Where("username LIKE ? OR email LIKE ?", searchPattern, searchPattern)

	// 获取总数
	if err := query.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}