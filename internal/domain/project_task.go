package domain

import (
	"time"

	"gorm.io/gorm"
)

// ProjectTask 项目任务记录
// 当用户接受项目时创建，用于跟踪用户的任务进度
type ProjectTask struct {
	ID        uint           `gorm:"primaryKey;comment:任务ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`

	// 关联项目
	ProjectID uint    `gorm:"not null;index;uniqueIndex:idx_task_project_user;comment:项目ID" json:"project_id"`
	Project   Project `gorm:"-" json:"project,omitempty"`

	// 接受任务的用户
	UserID uint `gorm:"not null;index;uniqueIndex:idx_task_project_user;comment:任务执行用户ID" json:"user_id"`
	User   User `gorm:"-" json:"user,omitempty"`

	// 任务状态: accepted(已接受/进行中)
	Status string `gorm:"size:20;default:'accepted';index;comment:任务状态(accepted)" json:"status"`

	// 接受任务时间
	AcceptedAt time.Time `gorm:"comment:接受任务时间" json:"accepted_at"`
}

// TableName 指定表名
func (ProjectTask) TableName() string {
	return "project_tasks"
}

// ProjectTaskRepository 项目任务仓库接口
type ProjectTaskRepository interface {
	Create(task *ProjectTask) error
	FindByID(id uint) (*ProjectTask, error)
	FindByProjectAndUser(projectID, userID uint) (*ProjectTask, error)
	FindByUserID(userID uint) ([]ProjectTask, error)
	FindAcceptedByUserID(userID uint) ([]ProjectTask, error) // 获取用户已接受的任务
	Update(task *ProjectTask) error
	Delete(id uint) error
}

// ProjectTaskService 项目任务服务接口
type ProjectTaskService interface {
	AcceptTask(projectID, userID uint) (*ProjectTask, error)
	GetUserTasks(userID uint) ([]ProjectTask, error)
	GetUserAcceptedProjectIDs(userID uint) ([]uint, error) // 获取用户已接受任务的项目ID列表
}
