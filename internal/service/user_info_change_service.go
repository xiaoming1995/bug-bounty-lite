package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type userInfoChangeService struct {
	repo domain.UserInfoChangeRepository
}

// NewUserInfoChangeService 创建用户信息变更服务实例
func NewUserInfoChangeService(repo domain.UserInfoChangeRepository) domain.UserInfoChangeService {
	return &userInfoChangeService{repo: repo}
}

// SubmitChangeRequest 提交用户信息变更申请
func (s *userInfoChangeService) SubmitChangeRequest(userID uint, phone, email, name string) (*domain.UserInfoChangeRequest, error) {
	// 1. 检查是否已有待审核的申请
	pending, err := s.repo.FindPendingByUserID(userID)
	if err != nil {
		return nil, err
	}
	if pending != nil {
		return nil, errors.New("您已有待审核的变更申请，请等待审核完成后再提交")
	}

	// 2. 创建新的变更申请
	request := &domain.UserInfoChangeRequest{
		UserID: userID,
		Phone:  phone,
		Email:  email,
		Name:   name,
		Status: "pending",
	}

	// 3. 保存到数据库
	if err := s.repo.Create(request); err != nil {
		return nil, err
	}

	return request, nil
}

// GetUserChangeRequests 获取用户的所有变更申请
func (s *userInfoChangeService) GetUserChangeRequests(userID uint) ([]domain.UserInfoChangeRequest, error) {
	return s.repo.FindByUserID(userID)
}

// GetChangeRequest 获取单个变更申请（只能查看自己的）
func (s *userInfoChangeService) GetChangeRequest(id uint, userID uint) (*domain.UserInfoChangeRequest, error) {
	request, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, errors.New("变更申请不存在")
	}

	// 只能查看自己的申请
	if request.UserID != userID {
		return nil, errors.New("无权访问此变更申请")
	}

	return request, nil
}

