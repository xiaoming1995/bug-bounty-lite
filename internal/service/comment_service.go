package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type commentService struct {
	repo       domain.CommentRepository
	reportRepo domain.ReportRepository
}

// NewCommentService 创建评论服务实例
func NewCommentService(repo domain.CommentRepository, reportRepo domain.ReportRepository) domain.CommentService {
	return &commentService{
		repo:       repo,
		reportRepo: reportRepo,
	}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(reportID uint, authorID uint, content string) (*domain.ReportComment, error) {
	// 验证内容不为空
	if content == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 验证报告存在
	_, err := s.reportRepo.FindByID(reportID)
	if err != nil {
		return nil, errors.New("漏洞报告不存在")
	}

	comment := &domain.ReportComment{
		ReportID: reportID,
		AuthorID: authorID,
		Content:  content,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetReportComments 获取报告的所有评论
func (s *commentService) GetReportComments(reportID uint) ([]domain.ReportComment, error) {
	return s.repo.FindByReportID(reportID)
}

// DeleteComment 删除评论（仅作者或管理员可删除）
func (s *commentService) DeleteComment(commentID uint, userID uint, userRole string) error {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// 权限校验：仅作者或管理员可删除
	if comment.AuthorID != userID && userRole != "admin" {
		return errors.New("无权删除此评论")
	}

	return s.repo.Delete(commentID)
}
