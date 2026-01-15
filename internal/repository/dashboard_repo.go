package repository

import (
	"bug-bounty-lite/internal/domain"
	"time"

	"gorm.io/gorm"
)

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepo(db *gorm.DB) domain.DashboardRepository {
	return &dashboardRepo{db: db}
}

// CountBySeverity 按危害等级统计已审核漏洞数量
// 只统计 status 不为 Pending 且 severity 不为空的报告
// authorID: nil=查询所有, 非nil=只查询指定用户的
func (r *dashboardRepo) CountBySeverity(authorID *uint) (*domain.SeverityStatistics, error) {
	stats := &domain.SeverityStatistics{}

	// 基础查询：已审核的报告（status 不为 Pending）
	baseQuery := r.db.Model(&domain.Report{}).Where("status != ?", "Pending")

	// 如果指定了用户ID，只查询该用户的报告
	if authorID != nil && *authorID > 0 {
		baseQuery = baseQuery.Where("author_id = ?", *authorID)
	}

	// 统计各等级数量
	var results []struct {
		Severity string
		Count    int64
	}

	err := baseQuery.Select("severity, COUNT(*) as count").
		Where("severity != '' AND severity IS NOT NULL").
		Group("severity").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 映射结果
	for _, item := range results {
		switch item.Severity {
		case "Critical":
			stats.Critical = item.Count
		case "High":
			stats.High = item.Count
		case "Medium":
			stats.Medium = item.Count
		case "Low":
			stats.Low = item.Count
		case "None":
			stats.None = item.Count
		}
	}

	// 计算总数
	stats.Total = stats.Critical + stats.High + stats.Medium + stats.Low + stats.None

	return stats, nil
}

// GetTrend 获取漏洞趋势数据
// authorID: nil=查询所有, 非nil=只查询指定用户的
func (r *dashboardRepo) GetTrend(period string, authorID *uint) ([]domain.TrendItem, error) {
	var results []domain.TrendItem

	// 基础查询：已审核的报告（status 不为 Pending），按 created_at 分组
	baseQuery := r.db.Model(&domain.Report{}).Where("status != ?", "Pending")

	// 如果指定了用户ID，只查询该用户的报告
	if authorID != nil && *authorID > 0 {
		baseQuery = baseQuery.Where("author_id = ?", *authorID)
	}

	now := time.Now()

	switch period {
	case "day":
		// 当月每天
		year, month, _ := now.Date()
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

		var dbResults []struct {
			Day   int
			Count int64
		}

		err := baseQuery.
			Where("created_at >= ? AND created_at <= ?", startOfMonth, endOfMonth).
			Select("DAY(created_at) as day, COUNT(*) as count").
			Group("DAY(created_at)").
			Order("day").
			Scan(&dbResults).Error

		if err != nil {
			return nil, err
		}

		// 填充当月所有天数
		daysInMonth := endOfMonth.Day()
		dayMap := make(map[int]int64)
		for _, item := range dbResults {
			dayMap[item.Day] = item.Count
		}

		for day := 1; day <= daysInMonth; day++ {
			results = append(results, domain.TrendItem{
				Label: time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Format("2日"),
				Value: dayMap[day],
			})
		}

	case "month":
		// 当年每月
		year := now.Year()
		startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
		endOfYear := time.Date(year, 12, 31, 23, 59, 59, 0, now.Location())

		var dbResults []struct {
			Month int
			Count int64
		}

		err := baseQuery.
			Where("created_at >= ? AND created_at <= ?", startOfYear, endOfYear).
			Select("MONTH(created_at) as month, COUNT(*) as count").
			Group("MONTH(created_at)").
			Order("month").
			Scan(&dbResults).Error

		if err != nil {
			return nil, err
		}

		// 填充所有月份
		monthLabels := []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
		monthMap := make(map[int]int64)
		for _, item := range dbResults {
			monthMap[item.Month] = item.Count
		}

		for i, label := range monthLabels {
			results = append(results, domain.TrendItem{
				Label: label,
				Value: monthMap[i+1],
			})
		}

	case "year":
		// 近5年
		currentYear := now.Year()
		startYear := currentYear - 4
		startOfRange := time.Date(startYear, 1, 1, 0, 0, 0, 0, now.Location())

		var dbResults []struct {
			Year  int
			Count int64
		}

		err := baseQuery.
			Where("created_at >= ?", startOfRange).
			Select("YEAR(created_at) as year, COUNT(*) as count").
			Group("YEAR(created_at)").
			Order("year").
			Scan(&dbResults).Error

		if err != nil {
			return nil, err
		}

		// 填充近5年
		yearMap := make(map[int]int64)
		for _, item := range dbResults {
			yearMap[item.Year] = item.Count
		}

		for year := startYear; year <= currentYear; year++ {
			results = append(results, domain.TrendItem{
				Label: time.Date(year, 1, 1, 0, 0, 0, 0, now.Location()).Format("2006"),
				Value: yearMap[year],
			})
		}

	default:
		return nil, nil
	}

	return results, nil
}

// ListByStatus 按状态获取漏洞列表
// authorID: nil=查询所有, 非nil=只查询指定用户的
func (r *dashboardRepo) ListByStatus(isPending bool, limit int, authorID *uint) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	query := r.db.Model(&domain.Report{})

	if isPending {
		query = query.Where("status = ?", "Pending")
	} else {
		query = query.Where("status != ?", "Pending")
	}

	// 如果指定了用户ID，只查询该用户的报告
	if authorID != nil && *authorID > 0 {
		query = query.Where("author_id = ?", *authorID)
	}

	// 获取总数
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	if err := query.Session(&gorm.Session{}).Order("created_at DESC").Limit(limit).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	// 加载关联数据
	for i := range reports {
		r.loadAssociations(&reports[i])
	}

	return reports, total, nil
}

// loadAssociations 手动加载关联数据
func (r *dashboardRepo) loadAssociations(report *domain.Report) {
	// 加载 Author
	if report.AuthorID > 0 {
		var author domain.User
		if err := r.db.First(&author, report.AuthorID).Error; err == nil {
			report.Author = author
		}
	}

	// 加载 Project
	if report.ProjectID > 0 {
		var project domain.Project
		if err := r.db.Unscoped().First(&project, report.ProjectID).Error; err == nil {
			report.Project = project
		}
	}

	// 加载 VulnerabilityType
	if report.VulnerabilityTypeID > 0 {
		var vulnType domain.SystemConfig
		if err := r.db.First(&vulnType, report.VulnerabilityTypeID).Error; err == nil {
			report.VulnerabilityType = vulnType
		}
	}

	// 加载 SelfAssessment
	if report.SelfAssessmentID != nil && *report.SelfAssessmentID > 0 {
		var selfAssessment domain.SystemConfig
		if err := r.db.First(&selfAssessment, *report.SelfAssessmentID).Error; err == nil {
			report.SelfAssessment = selfAssessment
		}
	}
}
