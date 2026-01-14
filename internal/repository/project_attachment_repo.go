package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type projectAttachmentRepo struct {
	db *gorm.DB
}

// NewProjectAttachmentRepository 创建项目附件仓库实例
func NewProjectAttachmentRepository(db *gorm.DB) domain.ProjectAttachmentRepository {
	return &projectAttachmentRepo{db: db}
}

// Create 创建附件记录
func (r *projectAttachmentRepo) Create(attachment *domain.ProjectAttachment) error {
	return r.db.Create(attachment).Error
}

// FindByProjectID 根据项目ID查找所有附件
func (r *projectAttachmentRepo) FindByProjectID(projectID uint) ([]domain.ProjectAttachment, error) {
	var attachments []domain.ProjectAttachment
	err := r.db.Where("project_id = ?", projectID).Order("sort_order ASC, id ASC").Find(&attachments).Error
	return attachments, err
}

// Delete 删除附件记录
func (r *projectAttachmentRepo) Delete(id uint) error {
	return r.db.Delete(&domain.ProjectAttachment{}, id).Error
}

// DeleteByProjectID 删除项目的所有附件
func (r *projectAttachmentRepo) DeleteByProjectID(projectID uint) error {
	return r.db.Where("project_id = ?", projectID).Delete(&domain.ProjectAttachment{}).Error
}
