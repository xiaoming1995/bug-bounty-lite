package repository

import (
	"bug-bounty-lite/internal/domain"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) domain.ProjectRepository {
	return &projectRepo{db: db}
}

// Create 创建项目
func (r *projectRepo) Create(project *domain.Project) error {
	return r.db.Create(project).Error
}

// FindByID 根据ID查找项目
func (r *projectRepo) FindByID(id uint) (*domain.Project, error) {
	var project domain.Project
	err := r.db.First(&project, id).Error
	return &project, err
}

// List 分页获取项目列表
func (r *projectRepo) List(page, pageSize int, includeInactive bool) ([]domain.Project, int64, error) {
	var projects []domain.Project
	var total int64

	// 计算 Offset
	offset := (page - 1) * pageSize

	// 开启查询会话
	query := r.db.Model(&domain.Project{})

	// 如果不包含非活跃项目，则过滤状态
	if !includeInactive {
		query = query.Where("status = ?", "active")
	}

	// 先查总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查当前页数据，按ID倒序（最新的在前）
	err := query.
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&projects).Error

	return projects, total, err
}

// Update 更新项目
func (r *projectRepo) Update(project *domain.Project) error {
	return r.db.Save(project).Error
}

// Delete 删除项目
func (r *projectRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Project{}, id).Error
}

