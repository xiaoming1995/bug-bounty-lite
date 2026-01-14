package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type articleLikeRepo struct {
	db *gorm.DB
}

// NewArticleLikeRepository 创建点赞仓库
func NewArticleLikeRepository(db *gorm.DB) domain.ArticleLikeRepository {
	return &articleLikeRepo{db: db}
}

// HasLiked 检查用户是否已点赞
func (r *articleLikeRepo) HasLiked(articleID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.ArticleLike{}).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Like 点赞
func (r *articleLikeRepo) Like(articleID, userID uint) error {
	like := &domain.ArticleLike{
		ArticleID: articleID,
		UserID:    userID,
	}
	return r.db.Create(like).Error
}

// Unlike 取消点赞
func (r *articleLikeRepo) Unlike(articleID, userID uint) error {
	return r.db.Where("article_id = ? AND user_id = ?", articleID, userID).
		Delete(&domain.ArticleLike{}).Error
}

// GetLikeCount 获取点赞数
func (r *articleLikeRepo) GetLikeCount(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.ArticleLike{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
