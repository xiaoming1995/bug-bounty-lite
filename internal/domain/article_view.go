package domain

import "time"

// ArticleView 文章访问记录（用于 IP 访问限制）
type ArticleView struct {
	ID        uint      `gorm:"primaryKey;comment:记录ID" json:"id"`
	ArticleID uint      `gorm:"index;not null;comment:文章ID" json:"article_id"`
	IP        string    `gorm:"size:45;index;not null;comment:访问IP" json:"ip"`
	ViewDate  string    `gorm:"size:10;index;not null;comment:访问日期(YYYY-MM-DD)" json:"view_date"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
}

func (ArticleView) TableName() string {
	return "article_views"
}

// ArticleViewRepository 文章访问记录仓库接口
type ArticleViewRepository interface {
	// HasViewedToday 检查指定 IP 今日是否已访问该文章
	HasViewedToday(articleID uint, ip string) (bool, error)
	// RecordView 记录访问
	RecordView(articleID uint, ip string) error
}
