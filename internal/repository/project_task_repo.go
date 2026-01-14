package repository

import (
	"bug-bounty-lite/internal/domain"
	"time"

	"gorm.io/gorm"
)

type projectTaskRepo struct {
	db *gorm.DB
}

// NewProjectTaskRepository 创建项目任务仓库实例
func NewProjectTaskRepository(db *gorm.DB) domain.ProjectTaskRepository {
	return &projectTaskRepo{db: db}
}

// Create 创建项目任务记录
func (r *projectTaskRepo) Create(task *domain.ProjectTask) error {
	if task.AcceptedAt.IsZero() {
		task.AcceptedAt = time.Now()
	}
	return r.db.Create(task).Error
}

// FindByID 根据ID查找任务
func (r *projectTaskRepo) FindByID(id uint) (*domain.ProjectTask, error) {
	var task domain.ProjectTask
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}

	// 手动加载关联的 Project
	if task.ProjectID > 0 {
		var project domain.Project
		if err := r.db.First(&project, task.ProjectID).Error; err == nil {
			task.Project = project
		}
	}

	return &task, nil
}

// FindByProjectAndUser 根据项目ID和用户ID查找任务
func (r *projectTaskRepo) FindByProjectAndUser(projectID, userID uint) (*domain.ProjectTask, error) {
	var task domain.ProjectTask
	err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// FindByUserID 根据用户ID查找所有任务
func (r *projectTaskRepo) FindByUserID(userID uint) ([]domain.ProjectTask, error) {
	var tasks []domain.ProjectTask
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	// 手动加载关联的 Project
	for i := range tasks {
		if tasks[i].ProjectID > 0 {
			var project domain.Project
			if err := r.db.First(&project, tasks[i].ProjectID).Error; err == nil {
				tasks[i].Project = project
			}
		}
	}

	return tasks, nil
}

// FindAcceptedByUserID 获取用户已接受的任务（状态为 accepted）
func (r *projectTaskRepo) FindAcceptedByUserID(userID uint) ([]domain.ProjectTask, error) {
	var tasks []domain.ProjectTask
	err := r.db.Where("user_id = ? AND status = ?", userID, "accepted").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	// 手动加载关联的 Project
	for i := range tasks {
		if tasks[i].ProjectID > 0 {
			var project domain.Project
			if err := r.db.First(&project, tasks[i].ProjectID).Error; err == nil {
				tasks[i].Project = project
			}
		}
	}

	return tasks, nil
}

// Update 更新任务
func (r *projectTaskRepo) Update(task *domain.ProjectTask) error {
	return r.db.Save(task).Error
}

// Delete 删除任务（软删除）
func (r *projectTaskRepo) Delete(id uint) error {
	return r.db.Delete(&domain.ProjectTask{}, id).Error
}
