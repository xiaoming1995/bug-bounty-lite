package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type projectAssignmentRepo struct {
	db *gorm.DB
}

// NewProjectAssignmentRepository 创建项目指派仓库实例
func NewProjectAssignmentRepository(db *gorm.DB) domain.ProjectAssignmentRepository {
	return &projectAssignmentRepo{db: db}
}

// Create 创建项目指派记录
func (r *projectAssignmentRepo) Create(assignment *domain.ProjectAssignment) error {
	return r.db.Create(assignment).Error
}

// FindByProjectAndUser 根据项目ID和用户ID查找指派记录
func (r *projectAssignmentRepo) FindByProjectAndUser(projectID, userID uint) (*domain.ProjectAssignment, error) {
	var assignment domain.ProjectAssignment
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// FindByProjectID 根据项目ID查找所有指派记录
func (r *projectAssignmentRepo) FindByProjectID(projectID uint) ([]domain.ProjectAssignment, error) {
	var assignments []domain.ProjectAssignment
	err := r.db.Where("project_id = ?", projectID).Find(&assignments).Error
	return assignments, err
}

// FindByUserID 根据用户ID查找所有指派给该用户的记录
func (r *projectAssignmentRepo) FindByUserID(userID uint) ([]domain.ProjectAssignment, error) {
	var assignments []domain.ProjectAssignment
	err := r.db.Where("user_id = ?", userID).Find(&assignments).Error
	return assignments, err
}

// Delete 删除指派记录
func (r *projectAssignmentRepo) Delete(id uint) error {
	return r.db.Delete(&domain.ProjectAssignment{}, id).Error
}

// DeleteByProjectAndUser 根据项目ID和用户ID删除指派记录
func (r *projectAssignmentRepo) DeleteByProjectAndUser(projectID, userID uint) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&domain.ProjectAssignment{}).Error
}
