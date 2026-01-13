package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type avatarRepo struct {
	db *gorm.DB
}

// NewAvatarRepo 创建头像仓库实例
func NewAvatarRepo(db *gorm.DB) domain.AvatarRepository {
	return &avatarRepo{db: db}
}

// Create 创建头像
func (r *avatarRepo) Create(avatar *domain.Avatar) error {
	return r.db.Create(avatar).Error
}

// FindByID 根据ID查找头像
func (r *avatarRepo) FindByID(id uint) (*domain.Avatar, error) {
	var avatar domain.Avatar
	if err := r.db.First(&avatar, id).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

// List 获取所有头像
func (r *avatarRepo) List() ([]domain.Avatar, error) {
	var avatars []domain.Avatar
	if err := r.db.Order("sort_order ASC, id ASC").Find(&avatars).Error; err != nil {
		return nil, err
	}
	return avatars, nil
}

// ListActive 获取启用的头像
func (r *avatarRepo) ListActive() ([]domain.Avatar, error) {
	var avatars []domain.Avatar
	if err := r.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&avatars).Error; err != nil {
		return nil, err
	}
	return avatars, nil
}

// Update 更新头像
func (r *avatarRepo) Update(avatar *domain.Avatar) error {
	return r.db.Save(avatar).Error
}

// Delete 删除头像
func (r *avatarRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Avatar{}, id).Error
}
