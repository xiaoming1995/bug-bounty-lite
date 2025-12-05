package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type systemConfigService struct {
	repo domain.SystemConfigRepository
}

func NewSystemConfigService(repo domain.SystemConfigRepository) domain.SystemConfigService {
	return &systemConfigService{repo: repo}
}

// GetConfigsByType 根据类型获取配置列表
func (s *systemConfigService) GetConfigsByType(configType string, includeInactive bool) ([]domain.SystemConfig, error) {
	if configType == "" {
		return nil, errors.New("配置类型不能为空")
	}
	return s.repo.FindByType(configType, includeInactive)
}

// GetConfig 获取配置详情
func (s *systemConfigService) GetConfig(id uint) (*domain.SystemConfig, error) {
	config, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("配置不存在")
	}
	return config, nil
}

// CreateConfig 创建配置
func (s *systemConfigService) CreateConfig(config *domain.SystemConfig) error {
	// 校验必填字段
	if config.ConfigType == "" {
		return errors.New("配置类型不能为空")
	}
	if config.ConfigKey == "" {
		return errors.New("配置键不能为空")
	}
	if config.ConfigValue == "" {
		return errors.New("配置值不能为空")
	}

	// 设置默认值
	if config.Status == "" {
		config.Status = "active"
	}

	return s.repo.Create(config)
}

// UpdateConfig 更新配置
func (s *systemConfigService) UpdateConfig(id uint, config *domain.SystemConfig) error {
	// 检查配置是否存在
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("配置不存在")
	}

	// 更新字段
	if config.ConfigType != "" {
		existing.ConfigType = config.ConfigType
	}
	if config.ConfigKey != "" {
		existing.ConfigKey = config.ConfigKey
	}
	if config.ConfigValue != "" {
		existing.ConfigValue = config.ConfigValue
	}
	if config.Description != "" {
		existing.Description = config.Description
	}
	if config.SortOrder != 0 || config.SortOrder != existing.SortOrder {
		existing.SortOrder = config.SortOrder
	}
	if config.Status != "" {
		existing.Status = config.Status
	}
	if config.ExtraData != nil {
		existing.ExtraData = config.ExtraData
	}

	return s.repo.Update(existing)
}

// DeleteConfig 删除配置
func (s *systemConfigService) DeleteConfig(id uint) error {
	// 检查配置是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("配置不存在")
	}

	return s.repo.Delete(id)
}

