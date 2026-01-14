package domain

import (
	"time"

	"gorm.io/gorm"
)

// Project 项目实体
type Project struct {
	ID        uint           `gorm:"primaryKey;comment:项目ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`

	// 项目名称
	Name string `gorm:"size:255;not null;comment:项目名称" json:"name"`

	// 项目描述
	Description string `gorm:"type:text;comment:项目描述" json:"description"`

	// 备注（平台内部使用）
	Note string `gorm:"type:text;comment:备注" json:"note"`

	// 项目难度: easy(简单), medium(中等), hard(困难), expert(专家)
	Difficulty string `gorm:"size:20;default:'medium';comment:项目难度(easy/medium/hard/expert)" json:"difficulty"`

	// 截止日期
	Deadline *time.Time `gorm:"comment:项目截止日期" json:"deadline"`

	// 项目状态: recruiting(招募中), in_progress(进行中), completed(已完成), closed(已关闭)
	Status string `gorm:"size:20;default:'recruiting';index;comment:项目状态(recruiting/in_progress/completed/closed)" json:"status"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// ProjectRepository 项目仓库接口
type ProjectRepository interface {
	Create(project *Project) error
	FindByID(id uint) (*Project, error)
	FindByIDWithDeleted(id uint) (*Project, error) // 包含已删除的项目
	List(page, pageSize int, includeInactive bool) ([]Project, int64, error)
	ListWithDeleted(page, pageSize int) ([]Project, int64, error) // 包含已删除的项目
	Update(project *Project) error
	Delete(id uint) error  // 软删除
	Restore(id uint) error // 恢复已删除的项目
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
	GetProjectWithDeleted(id uint) (*Project, error) // 包含已删除的项目
	ListProjects(page, pageSize int, includeInactive bool) ([]Project, int64, error)
	ListProjectsWithDeleted(page, pageSize int) ([]Project, int64, error) // 包含已删除的项目
	UpdateProject(id uint, input *ProjectUpdateInput) (*Project, error)
	DeleteProject(id uint) error  // 软删除
	RestoreProject(id uint) error // 恢复已删除的项目
}
