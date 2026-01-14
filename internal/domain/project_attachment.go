package domain

import (
	"time"
)

// ProjectAttachment 项目附件实体
type ProjectAttachment struct {
	ID        uint      `gorm:"primaryKey;comment:附件ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`

	// 关联项目
	ProjectID uint `gorm:"not null;index;comment:项目ID" json:"project_id"`

	// 附件名称
	Name string `gorm:"size:255;not null;comment:附件名称" json:"name"`

	// 附件URL
	URL string `gorm:"size:500;not null;comment:附件URL" json:"url"`

	// 文件大小（显示用，如 "256KB"）
	Size string `gorm:"size:50;comment:文件大小" json:"size"`

	// 文件类型（pdf, doc, image 等）
	Type string `gorm:"size:50;comment:文件类型" json:"type"`

	// 排序顺序
	SortOrder int `gorm:"default:0;comment:排序顺序" json:"sort_order"`
}

// TableName 指定表名
func (ProjectAttachment) TableName() string {
	return "project_attachments"
}

// ProjectAttachmentRepository 项目附件仓库接口
type ProjectAttachmentRepository interface {
	Create(attachment *ProjectAttachment) error
	FindByProjectID(projectID uint) ([]ProjectAttachment, error)
	Delete(id uint) error
	DeleteByProjectID(projectID uint) error
}
