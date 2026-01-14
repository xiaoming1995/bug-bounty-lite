package repository

import (
	"bug-bounty-lite/internal/domain"
	"time"

	"gorm.io/gorm"
)

type articleViewRepo struct {
	db *gorm.DB
}

// NewArticleViewRepo 创建文章访问记录仓库实例
func NewArticleViewRepo(db *gorm.DB) domain.ArticleViewRepository {
	return &articleViewRepo{db: db}
}

// HasViewedToday 检查指定 IP 今日是否已访问该文章
func (r *articleViewRepo) HasViewedToday(articleID uint, ip string) (bool, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := r.db.Model(&domain.ArticleView{}).
		Where("article_id = ? AND ip = ? AND view_date = ?", articleID, ip, today).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// RecordView 记录访问
func (r *articleViewRepo) RecordView(articleID uint, ip string) error {
	today := time.Now().Format("2006-01-02")
	view := &domain.ArticleView{
		ArticleID: articleID,
		IP:        ip,
		ViewDate:  today,
	}
	return r.db.Create(view).Error
}
