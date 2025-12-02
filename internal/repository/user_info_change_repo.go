package repository

import (
	"bug-bounty-lite/internal/domain"
	"errors"
	"gorm.io/gorm"
)

type userInfoChangeRepo struct {
	db *gorm.DB
}

// NewUserInfoChangeRepo 创建用户信息变更申请仓库实例
func NewUserInfoChangeRepo(db *gorm.DB) domain.UserInfoChangeRepository {
	return &userInfoChangeRepo{db: db}
}

// Create 创建变更申请
func (r *userInfoChangeRepo) Create(request *domain.UserInfoChangeRequest) error {
	return r.db.Create(request).Error
}

// FindByID 根据ID查找变更申请
func (r *userInfoChangeRepo) FindByID(id uint) (*domain.UserInfoChangeRequest, error) {
	var request domain.UserInfoChangeRequest
	err := r.db.First(&request, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// FindByUserID 查找用户的所有变更申请
func (r *userInfoChangeRepo) FindByUserID(userID uint) ([]domain.UserInfoChangeRequest, error) {
	var requests []domain.UserInfoChangeRequest
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

// FindPendingByUserID 查找用户待审核的变更申请
func (r *userInfoChangeRepo) FindPendingByUserID(userID uint) (*domain.UserInfoChangeRequest, error) {
	var request domain.UserInfoChangeRequest
	err := r.db.Where("user_id = ? AND status = ?", userID, "pending").
		First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// Update 更新变更申请
func (r *userInfoChangeRepo) Update(request *domain.UserInfoChangeRequest) error {
	return r.db.Save(request).Error
}

