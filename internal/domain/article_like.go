package domain

import "time"

// ArticleLike 文章点赞记录
type ArticleLike struct {
	ID        uint      `gorm:"primaryKey;comment:记录ID" json:"id"`
	ArticleID uint      `gorm:"index:idx_article_user,unique;not null;comment:文章ID" json:"article_id"`
	UserID    uint      `gorm:"index:idx_article_user,unique;not null;comment:用户ID" json:"user_id"`
	CreatedAt time.Time `gorm:"comment:点赞时间" json:"created_at"`
}

func (ArticleLike) TableName() string {
	return "article_likes"
}

// ArticleLikeRepository 点赞仓库接口
type ArticleLikeRepository interface {
	HasLiked(articleID, userID uint) (bool, error)
	Like(articleID, userID uint) error
	Unlike(articleID, userID uint) error
	GetLikeCount(articleID uint) (int64, error)
}
