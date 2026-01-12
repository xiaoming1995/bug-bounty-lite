package repository

import (
	"bug-bounty-lite/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 就是编译器找不到的那个函数！
// 它负责创建 userRepo 实例并返回接口
func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

// Create 创建用户
func (r *userRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByUsername 根据用户名查找
func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	// First 查不到数据会返回 ErrRecordNotFound 错误
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到不算系统错误，返回 nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找
func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Org").First(&user, id).Error
	return &user, err
}

// Update 更新用户
func (r *userRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
