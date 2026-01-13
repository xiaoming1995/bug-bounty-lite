package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type articleService struct {
	repo domain.ArticleRepository
}

// NewArticleService 创建文章服务实例
func NewArticleService(repo domain.ArticleRepository) domain.ArticleService {
	return &articleService{repo: repo}
}

// CreateArticle 创建文章
func (s *articleService) CreateArticle(authorID uint, title, description, content string) (*domain.Article, error) {
	if title == "" {
		return nil, errors.New("文章标题不能为空")
	}
	if content == "" {
		return nil, errors.New("文章内容不能为空")
	}

	article := &domain.Article{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		Content:     content,
		Status:      "pending", // 默认待审核
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

// UpdateArticle 更新文章
func (s *articleService) UpdateArticle(articleID, userID uint, title, description, content string) (*domain.Article, error) {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 权限校验：仅作者可编辑
	if article.AuthorID != userID {
		return nil, errors.New("无权编辑此文章")
	}

	// 状态校验：已发布的文章不能编辑
	if article.Status == "approved" {
		return nil, errors.New("已发布的文章不能编辑")
	}

	// 更新字段
	article.Title = title
	article.Description = description
	article.Content = content
	article.Status = "pending" // 重新提交后变为待审核

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}

// DeleteArticle 删除文章
func (s *articleService) DeleteArticle(articleID, userID uint, userRole string) error {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return errors.New("文章不存在")
	}

	// 权限校验：仅作者或管理员可删除
	if article.AuthorID != userID && userRole != "admin" {
		return errors.New("无权删除此文章")
	}

	// 状态校验：已发布的文章只有管理员可删除
	if article.Status == "approved" && userRole != "admin" {
		return errors.New("已发布的文章不能删除")
	}

	return s.repo.Delete(articleID)
}

// GetArticle 获取文章详情
func (s *articleService) GetArticle(id uint, incrementView bool) (*domain.Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 增加浏览量
	if incrementView && article.Status == "approved" {
		_ = s.repo.IncrementViews(id)
		article.Views++
	}

	return article, nil
}

// GetMyArticles 获取用户的文章列表
func (s *articleService) GetMyArticles(authorID uint) ([]domain.Article, error) {
	return s.repo.FindByAuthorID(authorID)
}

// GetPublishedArticles 获取已发布的文章列表（学习中心）
func (s *articleService) GetPublishedArticles() ([]domain.Article, error) {
	return s.repo.FindPublished()
}

// ReviewArticle 审核文章（管理员）
func (s *articleService) ReviewArticle(articleID uint, approved bool, rejectReason string) (*domain.Article, error) {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	if approved {
		article.Status = "approved"
		article.RejectReason = ""
	} else {
		article.Status = "rejected"
		article.RejectReason = rejectReason
	}

	if err := s.repo.Update(article); err != nil {
		return nil, err
	}

	return article, nil
}
