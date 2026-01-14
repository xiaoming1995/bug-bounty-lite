package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

// ArticleLikeCommentService 文章点赞评论服务接口
type ArticleLikeCommentService interface {
	ToggleLike(articleID, userID uint) (liked bool, likeCount int64, err error)
	GetLikeStatus(articleID, userID uint) (liked bool, likeCount int64, err error)
	AddComment(articleID, userID uint, content string) (*domain.ArticleComment, error)
	GetComments(articleID uint) ([]domain.ArticleComment, error)
	DeleteComment(commentID, userID uint) error
}

type articleLikeCommentService struct {
	likeRepo    domain.ArticleLikeRepository
	commentRepo domain.ArticleCommentRepository
	articleRepo domain.ArticleRepository
}

// NewArticleLikeCommentService 创建点赞评论服务
func NewArticleLikeCommentService(
	likeRepo domain.ArticleLikeRepository,
	commentRepo domain.ArticleCommentRepository,
	articleRepo domain.ArticleRepository,
) ArticleLikeCommentService {
	return &articleLikeCommentService{
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

// ToggleLike 切换点赞状态
func (s *articleLikeCommentService) ToggleLike(articleID, userID uint) (bool, int64, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return false, 0, errors.New("文章不存在")
	}
	if article.Status != "approved" {
		return false, 0, errors.New("文章未发布")
	}

	// 检查是否已点赞
	hasLiked, err := s.likeRepo.HasLiked(articleID, userID)
	if err != nil {
		return false, 0, err
	}

	if hasLiked {
		// 取消点赞
		if err := s.likeRepo.Unlike(articleID, userID); err != nil {
			return false, 0, err
		}
	} else {
		// 点赞
		if err := s.likeRepo.Like(articleID, userID); err != nil {
			return false, 0, err
		}
	}

	// 获取最新点赞数
	count, err := s.likeRepo.GetLikeCount(articleID)
	if err != nil {
		return !hasLiked, 0, err
	}

	// 同步更新文章的点赞数
	_ = s.articleRepo.UpdateLikes(articleID, int(count))

	return !hasLiked, count, nil
}

// GetLikeStatus 获取点赞状态
func (s *articleLikeCommentService) GetLikeStatus(articleID, userID uint) (bool, int64, error) {
	hasLiked, err := s.likeRepo.HasLiked(articleID, userID)
	if err != nil {
		return false, 0, err
	}

	count, err := s.likeRepo.GetLikeCount(articleID)
	if err != nil {
		return hasLiked, 0, err
	}

	return hasLiked, count, nil
}

// AddComment 发表评论
func (s *articleLikeCommentService) AddComment(articleID, userID uint, content string) (*domain.ArticleComment, error) {
	if content == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	if article.Status != "approved" {
		return nil, errors.New("文章未发布")
	}

	comment := &domain.ArticleComment{
		ArticleID: articleID,
		UserID:    userID,
		Content:   content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetComments 获取评论列表
func (s *articleLikeCommentService) GetComments(articleID uint) ([]domain.ArticleComment, error) {
	return s.commentRepo.FindByArticleID(articleID)
}

// DeleteComment 删除评论
func (s *articleLikeCommentService) DeleteComment(commentID, userID uint) error {
	return s.commentRepo.Delete(commentID, userID)
}
