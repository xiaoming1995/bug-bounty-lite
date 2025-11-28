package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type reportService struct {
	repo domain.ReportRepository
}

func NewReportService(repo domain.ReportRepository) domain.ReportService {
	return &reportService{repo: repo}
}

// SubmitReport 提交漏洞
func (s *reportService) SubmitReport(report *domain.Report) error {
	// 1. 强制初始化状态
	// 不管前端传什么状态，后端强制设为 Pending
	report.Status = "Pending"

	// 2. 校验参数
	if report.Title == "" {
		return errors.New("title is required")
	}
	if report.AuthorID == 0 {
		return errors.New("author is required")
	}

	// 3. 调用 Repo
	return s.repo.Create(report)
}

func (s *reportService) GetReport(id uint) (*domain.Report, error) {
	return s.repo.FindByID(id)
}

func (s *reportService) ListReports(page, pageSize int) ([]domain.Report, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}