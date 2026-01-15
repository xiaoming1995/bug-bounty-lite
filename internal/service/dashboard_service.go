package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"
)

type dashboardService struct {
	repo domain.DashboardRepository
}

func NewDashboardService(repo domain.DashboardRepository) domain.DashboardService {
	return &dashboardService{repo: repo}
}

// GetStatistics 获取漏洞统计数据
// whitehat 只能查看自己的统计，admin/vendor 可以查看所有
func (s *dashboardService) GetStatistics(userID uint, userRole string) (*domain.SeverityStatistics, error) {
	var authorID *uint
	if userRole == "whitehat" {
		authorID = &userID
	}
	return s.repo.CountBySeverity(authorID)
}

// GetTrend 获取漏洞趋势数据
// whitehat 只能查看自己的趋势，admin/vendor 可以查看所有
func (s *dashboardService) GetTrend(period string, userID uint, userRole string) ([]domain.TrendItem, error) {
	// 校验 period 参数
	if period != "day" && period != "month" && period != "year" {
		return nil, errors.New("invalid period, must be 'day', 'month' or 'year'")
	}

	var authorID *uint
	if userRole == "whitehat" {
		authorID = &userID
	}
	return s.repo.GetTrend(period, authorID)
}

// GetReportsByType 按类型获取漏洞列表
// whitehat 只能查看自己的报告，admin/vendor 可以查看所有
func (s *dashboardService) GetReportsByType(reportType string, limit int, userID uint, userRole string) ([]domain.Report, int64, error) {
	// 校验 reportType 参数
	if reportType != "pending" && reportType != "reviewed" {
		return nil, 0, errors.New("invalid report type, must be 'pending' or 'reviewed'")
	}

	// 校验 limit 参数
	if limit <= 0 || limit > 100 {
		limit = 6 // 默认值
	}

	var authorID *uint
	if userRole == "whitehat" {
		authorID = &userID
	}

	isPending := reportType == "pending"
	return s.repo.ListByStatus(isPending, limit, authorID)
}
