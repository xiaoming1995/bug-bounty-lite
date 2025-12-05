package domain

import (
	"time"
)

// Project 项目实体
type Project struct {
	ID        uint      `gorm:"primaryKey;comment:项目ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 项目名称
	Name string `gorm:"size:255;not null;comment:项目名称" json:"name"`

	// 项目描述
	Description string `gorm:"type:text;comment:项目描述" json:"description"`

	// 备注（平台内部使用）
	Note string `gorm:"type:text;comment:备注" json:"note"`

	// 项目状态: active(活跃), inactive(非活跃)
	Status string `gorm:"size:20;default:'active';index;comment:项目状态(active/inactive)" json:"status"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// ProjectRepository 项目仓库接口
type ProjectRepository interface {
	Create(project *Project) error
	FindByID(id uint) (*Project, error)
	List(page, pageSize int, includeInactive bool) ([]Project, int64, error)
	Update(project *Project) error
	Delete(id uint) error
}

// ProjectUpdateInput 更新项目输入
type ProjectUpdateInput struct {
	Name        string
	Description string
	Note        string
	Status      string
}

// ProjectService 项目服务接口
type ProjectService interface {
	CreateProject(project *Project) error
	GetProject(id uint, includeInactive bool) (*Project, error)
	ListProjects(page, pageSize int, includeInactive bool) ([]Project, int64, error)
	UpdateProject(id uint, input *ProjectUpdateInput) (*Project, error)
	DeleteProject(id uint) error
}

