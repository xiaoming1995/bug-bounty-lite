package domain

import "time"

// ArticleComment 文章评论
type ArticleComment struct {
	ID        uint      `gorm:"primaryKey;comment:评论ID" json:"id"`
	ArticleID uint      `gorm:"index;not null;comment:文章ID" json:"article_id"`
	UserID    uint      `gorm:"index;not null;comment:用户ID" json:"user_id"`
	Content   string    `gorm:"type:text;not null;comment:评论内容" json:"content"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 关联用户（手动加载）
	User *User `gorm:"-" json:"user,omitempty"`
}

func (ArticleComment) TableName() string {
	return "article_comments"
}

// ArticleCommentRepository 评论仓库接口
type ArticleCommentRepository interface {
	Create(comment *ArticleComment) error
	FindByArticleID(articleID uint) ([]ArticleComment, error)
	Delete(id, userID uint) error
	CountByArticleID(articleID uint) (int64, error)
}
