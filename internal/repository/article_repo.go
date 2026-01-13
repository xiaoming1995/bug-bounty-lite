package repository

import (
	"bug-bounty-lite/internal/domain"

	"gorm.io/gorm"
)

type articleRepo struct {
	db *gorm.DB
}

// NewArticleRepo 创建文章仓库实例
func NewArticleRepo(db *gorm.DB) domain.ArticleRepository {
	return &articleRepo{db: db}
}

// Create 创建文章
func (r *articleRepo) Create(article *domain.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *articleRepo) Update(article *domain.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *articleRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}

// FindByID 根据ID获取文章
func (r *articleRepo) FindByID(id uint) (*domain.Article, error) {
	var article domain.Article
	if err := r.db.First(&article, id).Error; err != nil {
		return nil, err
	}

	// 手动加载作者信息
	if article.AuthorID > 0 {
		var user domain.User
		if err := r.db.First(&user, article.AuthorID).Error; err == nil {
			// 加载头像
			if user.AvatarID > 0 {
				var avatar domain.Avatar
				if err := r.db.First(&avatar, user.AvatarID).Error; err == nil {
					user.Avatar = &avatar
				}
			}
			article.Author = &user
		}
	}

	return &article, nil
}

// FindByAuthorID 根据作者ID获取文章列表
func (r *articleRepo) FindByAuthorID(authorID uint) ([]domain.Article, error) {
	var articles []domain.Article
	if err := r.db.Where("author_id = ?", authorID).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// FindPublished 获取所有已发布的文章
func (r *articleRepo) FindPublished() ([]domain.Article, error) {
	var articles []domain.Article
	if err := r.db.Where("status = ?", "approved").Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, err
	}

	// 手动加载作者信息
	for i := range articles {
		if articles[i].AuthorID > 0 {
			var user domain.User
			if err := r.db.First(&user, articles[i].AuthorID).Error; err == nil {
				if user.AvatarID > 0 {
					var avatar domain.Avatar
					if err := r.db.First(&avatar, user.AvatarID).Error; err == nil {
						user.Avatar = &avatar
					}
				}
				articles[i].Author = &user
			}
		}
	}

	return articles, nil
}

// IncrementViews 增加浏览量
func (r *articleRepo) IncrementViews(id uint) error {
	return r.db.Model(&domain.Article{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}
