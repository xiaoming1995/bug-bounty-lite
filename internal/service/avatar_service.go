package service

import (
	"bug-bounty-lite/internal/domain"
)

type avatarService struct {
	repo domain.AvatarRepository
}

// NewAvatarService 创建头像服务实例
func NewAvatarService(repo domain.AvatarRepository) domain.AvatarService {
	return &avatarService{repo: repo}
}

// UploadAvatar 上传头像（记录到数据库）
func (s *avatarService) UploadAvatar(name string, url string) (*domain.Avatar, error) {
	avatar := &domain.Avatar{
		Name:     name,
		URL:      url,
		IsActive: true,
	}
	if err := s.repo.Create(avatar); err != nil {
		return nil, err
	}
	return avatar, nil
}

// GetAvatar 获取单个头像
func (s *avatarService) GetAvatar(id uint) (*domain.Avatar, error) {
	return s.repo.FindByID(id)
}

// ListAvatars 获取所有头像
func (s *avatarService) ListAvatars() ([]domain.Avatar, error) {
	return s.repo.List()
}

// ListActiveAvatars 获取启用的头像
func (s *avatarService) ListActiveAvatars() ([]domain.Avatar, error) {
	return s.repo.ListActive()
}

// UpdateAvatar 更新头像信息
func (s *avatarService) UpdateAvatar(id uint, name string, isActive bool, sortOrder int) (*domain.Avatar, error) {
	avatar, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	avatar.Name = name
	avatar.IsActive = isActive
	avatar.SortOrder = sortOrder
	if err := s.repo.Update(avatar); err != nil {
		return nil, err
	}
	return avatar, nil
}

// DeleteAvatar 删除头像
func (s *avatarService) DeleteAvatar(id uint) error {
	return s.repo.Delete(id)
}
