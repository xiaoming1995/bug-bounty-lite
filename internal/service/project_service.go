package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type projectService struct {
	repo domain.ProjectRepository
}

func NewProjectService(repo domain.ProjectRepository) domain.ProjectService {
	return &projectService{repo: repo}
}

// CreateProject 创建项目
func (s *projectService) CreateProject(project *domain.Project) error {
	// 1. 校验必填字段
	if project.Name == "" {
		return errors.New("项目名称不能为空")
	}

	// 2. 设置默认状态
	if project.Status == "" {
		project.Status = "active"
	}

	// 3. 调用 Repository
	return s.repo.Create(project)
}

// GetProject 获取项目详情（不包含已删除的）
func (s *projectService) GetProject(id uint, includeInactive bool) (*domain.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 如果不包含非活跃项目，且项目状态为非活跃，则返回错误
	if !includeInactive && project.Status != "active" {
		return nil, errors.New("项目不存在")
	}

	return project, nil
}

// GetProjectWithDeleted 获取项目详情（包含已删除的，仅管理员使用）
func (s *projectService) GetProjectWithDeleted(id uint) (*domain.Project, error) {
	project, err := s.repo.FindByIDWithDeleted(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}
	return project, nil
}

// ListProjects 获取项目列表（不包含已删除的）
func (s *projectService) ListProjects(page, pageSize int, includeInactive bool) ([]domain.Project, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(page, pageSize, includeInactive)
}

// ListProjectsWithDeleted 获取项目列表（包含已删除的，仅管理员使用）
func (s *projectService) ListProjectsWithDeleted(page, pageSize int) ([]domain.Project, int64, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.ListWithDeleted(page, pageSize)
}

// UpdateProject 更新项目
func (s *projectService) UpdateProject(id uint, input *domain.ProjectUpdateInput) (*domain.Project, error) {
	// 1. 获取现有项目
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 2. 更新字段（只更新提供的非空字段）
	if input.Name != "" {
		project.Name = input.Name
	}
	if input.Description != "" {
		project.Description = input.Description
	}
	if input.Note != "" {
		project.Note = input.Note
	}
	if input.Status != "" {
		// 校验状态值
		if input.Status != "active" && input.Status != "inactive" {
			return nil, errors.New("无效的项目状态")
		}
		project.Status = input.Status
	}

	// 3. 保存
	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject 软删除项目
func (s *projectService) DeleteProject(id uint) error {
	// 检查项目是否存在
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("项目不存在")
	}

	return s.repo.Delete(id)
}

// RestoreProject 恢复已删除的项目
func (s *projectService) RestoreProject(id uint) error {
	// 检查项目是否存在（包含已删除的）
	project, err := s.repo.FindByIDWithDeleted(id)
	if err != nil {
		return errors.New("项目不存在")
	}

	// 检查项目是否已被删除
	if !project.DeletedAt.Valid {
		return errors.New("项目未被删除，无需恢复")
	}

	return s.repo.Restore(id)
}
