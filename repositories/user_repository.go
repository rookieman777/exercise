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
	// FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	// FindByUsername(username string) (*models.User, error)
	// Update(user *models.User) error
	// Delete(id uint) error
	// FindAll(page, pageSize int) ([]models.User, int64, error)
	// Search(keyword string, page, pageSize int) ([]models.User, int64, error)
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

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
