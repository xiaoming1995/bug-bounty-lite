package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type articleCommentRepo struct {
	db *gorm.DB
}

// NewArticleCommentRepository 创建评论仓库
func NewArticleCommentRepository(db *gorm.DB) domain.ArticleCommentRepository {
	return &articleCommentRepo{db: db}
}

// Create 创建评论
func (r *articleCommentRepo) Create(comment *domain.ArticleComment) error {
	return r.db.Create(comment).Error
}

// FindByArticleID 获取文章评论列表（按时间倒序）
func (r *articleCommentRepo) FindByArticleID(articleID uint) ([]domain.ArticleComment, error) {
	var comments []domain.ArticleComment
	err := r.db.Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// 手动加载用户信息
	for i := range comments {
		var user domain.User
		if err := r.db.First(&user, comments[i].UserID).Error; err == nil {
			// 加载头像
			if user.AvatarID > 0 {
				var avatar domain.Avatar
				if err := r.db.First(&avatar, user.AvatarID).Error; err == nil {
					user.Avatar = &avatar
				}
			}
			comments[i].User = &user
		}
	}

	return comments, nil
}

// Delete 删除评论（仅评论者可删除）
func (r *articleCommentRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.ArticleComment{}).Error
}

// CountByArticleID 获取文章评论数
func (r *articleCommentRepo) CountByArticleID(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.ArticleComment{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
