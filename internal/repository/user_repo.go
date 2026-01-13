package repository

import (
	"bug-bounty-lite/internal/domain"
	"errors"
	"time"

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
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	// 手动加载组织信息（代码逻辑关联，非外键）
	if user.OrgID > 0 {
		var org domain.Organization
		if err := r.db.First(&org, user.OrgID).Error; err == nil {
			user.Org = &org
		}
	}

	// 手动加载头像信息
	if user.AvatarID > 0 {
		var avatar domain.Avatar
		if err := r.db.First(&avatar, user.AvatarID).Error; err == nil {
			user.Avatar = &avatar
		}
	}

	return &user, nil
}

func (r *userRepo) Update(user *domain.User) error {
	// 使用 Save 进行全量更新
	return r.db.Save(user).Error
}

// UpdateLastLoginAt 专门更新最后登录时间
func (r *userRepo) UpdateLastLoginAt(userID uint, loginTime time.Time) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("last_login_at", loginTime).Error
}

// UpdateProfileFields 仅更新个人资料字段（避免外键约束问题）
func (r *userRepo) UpdateProfileFields(userID uint, name, bio, phone, email string) error {
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if bio != "" {
		updates["bio"] = bio
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if email != "" {
		updates["email"] = email
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateAvatarID 更新用户头像
func (r *userRepo) UpdateAvatarID(userID uint, avatarID uint) error {
	return r.db.Model(&domain.User{}).Where("id = ?", userID).Update("avatar_id", avatarID).Error
}
