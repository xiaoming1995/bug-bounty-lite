package domain

import "time"

// Article 文章实体
type Article struct {
	ID        uint      `gorm:"primaryKey;comment:文章ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 文章标题
	Title string `gorm:"size:200;not null;comment:文章标题" json:"title"`

	// 简要描述
	Description string `gorm:"size:500;comment:简要描述" json:"description"`

	// 文章内容 (HTML)
	Content string `gorm:"type:longtext;comment:文章内容(HTML)" json:"content"`

	// 作者ID
	AuthorID uint  `gorm:"index;not null;comment:作者ID" json:"author_id"`
	Author   *User `gorm:"-" json:"author,omitempty"` // 手动加载

	// 状态: pending(待审核), approved(已发布), rejected(被驳回)
	Status string `gorm:"size:20;default:'pending';index;comment:状态(pending:待审核, approved:已发布, rejected:驳回)" json:"status"`

	// 驳回原因
	RejectReason string `gorm:"size:500;comment:驳回原因" json:"reject_reason,omitempty"`

	// 分类
	Category string `gorm:"size:50;index;comment:文章分类" json:"category"`

	// 是否精选
	IsFeatured bool `gorm:"default:false;index;comment:是否精选" json:"is_featured"`

	// 统计数据
	Views int `gorm:"default:0;comment:浏览量" json:"views"`
	Likes int `gorm:"default:0;comment:点赞量" json:"likes"`
}

func (Article) TableName() string {
	return "articles"
}

// ArticleRepository 文章仓库接口
type ArticleRepository interface {
	Create(article *Article) error
	Update(article *Article) error
	Delete(id uint) error
	FindByID(id uint) (*Article, error)
	FindByAuthorID(authorID uint) ([]Article, error)
	FindPublished() ([]Article, error)
	FindFeatured(limit int) ([]Article, error)
	FindHot(limit int) ([]Article, error)
	SetFeatured(id uint, featured bool) error
	IncrementViews(id uint) error
	UpdateLikes(id uint, likes int) error // 更新点赞数
}

// ArticleService 文章服务接口
type ArticleService interface {
	CreateArticle(authorID uint, userRole, title, description, content, category string) (*Article, error)
	UpdateArticle(articleID, userID uint, title, description, content, category string) (*Article, error)
	DeleteArticle(articleID, userID uint, userRole string) error
	GetArticle(id uint, incrementView bool, clientIP string) (*Article, error)
	GetMyArticles(authorID uint) ([]Article, error)
	GetPublishedArticles() ([]Article, error)
	GetFeaturedArticles(limit int) ([]Article, error)
	GetHotArticles(limit int) ([]Article, error)
	SetFeatured(articleID uint, featured bool) error
	ReviewArticle(articleID uint, approved bool, rejectReason string) (*Article, error)
}
