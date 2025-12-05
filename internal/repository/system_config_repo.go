package repository

import (
	"bug-bounty-lite/internal/domain"
	"gorm.io/gorm"
)

type systemConfigRepo struct {
	db *gorm.DB
}

func NewSystemConfigRepo(db *gorm.DB) domain.SystemConfigRepository {
	return &systemConfigRepo{db: db}
}

// Create 创建配置
func (r *systemConfigRepo) Create(config *domain.SystemConfig) error {
	return r.db.Create(config).Error
}

// FindByID 根据ID查找配置
func (r *systemConfigRepo) FindByID(id uint) (*domain.SystemConfig, error) {
	var config domain.SystemConfig
	err := r.db.First(&config, id).Error
	return &config, err
}

// FindByType 根据类型查找配置列表
func (r *systemConfigRepo) FindByType(configType string, includeInactive bool) ([]domain.SystemConfig, error) {
	var configs []domain.SystemConfig
	query := r.db.Where("config_type = ?", configType)

	if !includeInactive {
		query = query.Where("status = ?", "active")
	}

	err := query.Order("sort_order ASC, id ASC").Find(&configs).Error
	return configs, err
}

// Update 更新配置
func (r *systemConfigRepo) Update(config *domain.SystemConfig) error {
	return r.db.Save(config).Error
}

// Delete 删除配置
func (r *systemConfigRepo) Delete(id uint) error {
	return r.db.Delete(&domain.SystemConfig{}, id).Error
}

