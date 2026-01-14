package domain

import (
	"time"
)

// ProjectAssignment 项目指派记录
// 记录管理员将项目指派给哪些用户，只有被指派的用户才能在项目大厅看到该项目
type ProjectAssignment struct {
	ID        uint      `gorm:"primaryKey;comment:指派ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:指派时间" json:"created_at"`

	// 关联项目
	ProjectID uint    `gorm:"not null;index;uniqueIndex:idx_project_user;comment:项目ID" json:"project_id"`
	Project   Project `gorm:"-" json:"project,omitempty"`

	// 被指派的用户
	UserID uint `gorm:"not null;index;uniqueIndex:idx_project_user;comment:被指派用户ID" json:"user_id"`
	User   User `gorm:"-" json:"user,omitempty"`
}

// TableName 指定表名
func (ProjectAssignment) TableName() string {
	return "project_assignments"
}

// ProjectAssignmentRepository 项目指派仓库接口
type ProjectAssignmentRepository interface {
	Create(assignment *ProjectAssignment) error
	FindByProjectAndUser(projectID, userID uint) (*ProjectAssignment, error)
	FindByProjectID(projectID uint) ([]ProjectAssignment, error)
	FindByUserID(userID uint) ([]ProjectAssignment, error)
	Delete(id uint) error
	DeleteByProjectAndUser(projectID, userID uint) error
}
