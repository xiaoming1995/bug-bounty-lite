package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepo 创建评论仓库实例
func NewCommentRepo(db *gorm.DB) domain.CommentRepository {
	return &commentRepo{db: db}
}

// Create 创建评论
func (r *commentRepo) Create(comment *domain.ReportComment) error {
	return r.db.Create(comment).Error
}

// FindByReportID 根据报告ID获取所有评论
func (r *commentRepo) FindByReportID(reportID uint) ([]domain.ReportComment, error) {
	var comments []domain.ReportComment
	if err := r.db.Where("report_id = ?", reportID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}

	// 手动加载作者信息
	for i := range comments {
		if comments[i].AuthorID > 0 {
			var user domain.User
			if err := r.db.First(&user, comments[i].AuthorID).Error; err == nil {
				// 加载头像信息
				if user.AvatarID > 0 {
					var avatar domain.Avatar
					if err := r.db.First(&avatar, user.AvatarID).Error; err == nil {
						user.Avatar = &avatar
					}
				}
				comments[i].Author = &user
			}
		}
	}

	return comments, nil
}

// FindByID 根据ID获取评论
func (r *commentRepo) FindByID(id uint) (*domain.ReportComment, error) {
	var comment domain.ReportComment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Delete 删除评论
func (r *commentRepo) Delete(id uint) error {
	return r.db.Delete(&domain.ReportComment{}, id).Error
}
