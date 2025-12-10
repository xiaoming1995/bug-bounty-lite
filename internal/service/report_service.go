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

	// 2. 校验必填参数
	if report.ProjectID == 0 {
		return errors.New("项目ID不能为空")
	}
	if report.VulnerabilityName == "" {
		return errors.New("漏洞名称不能为空")
	}
	if report.VulnerabilityTypeID == 0 {
		return errors.New("漏洞类型不能为空")
	}
	if report.AuthorID == 0 {
		return errors.New("提交者ID不能为空")
	}

	// 3. 设置默认值
	if report.Severity == "" {
		report.Severity = "Low"
	}

	// 4. 调用 Repo 创建
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

// UpdateReport 更新报告
func (s *reportService) UpdateReport(id uint, userID uint, userRole string, input *domain.ReportUpdateInput) (*domain.Report, error) {
	// 1. 获取现有报告
	report, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("report not found")
	}

	// 2. 权限校验
	// 只有报告作者、管理员或厂商可以更新
	if report.AuthorID != userID && userRole != "admin" && userRole != "vendor" {
		return nil, errors.New("permission denied")
	}

	// 3. 状态更新权限校验
	// 只有管理员和厂商可以修改状态
	if input.Status != "" && input.Status != report.Status {
		if userRole != "admin" && userRole != "vendor" {
			return nil, errors.New("only admin or vendor can change status")
		}
		// 校验状态流转
		if !isValidStatusTransition(report.Status, input.Status) {
			return nil, errors.New("invalid status transition")
		}
		report.Status = input.Status
	}

	// 4. 更新字段
	if input.ProjectID != 0 {
		report.ProjectID = input.ProjectID
	}
	if input.VulnerabilityName != "" {
		report.VulnerabilityName = input.VulnerabilityName
	}
	if input.VulnerabilityTypeID != 0 {
		report.VulnerabilityTypeID = input.VulnerabilityTypeID
	}
	if input.VulnerabilityImpact != "" {
		report.VulnerabilityImpact = input.VulnerabilityImpact
	}
	if input.SelfAssessment != "" {
		report.SelfAssessment = input.SelfAssessment
	}
	if input.VulnerabilityURL != "" {
		report.VulnerabilityURL = input.VulnerabilityURL
	}
	if input.VulnerabilityDetail != "" {
		report.VulnerabilityDetail = input.VulnerabilityDetail
	}
	if input.AttachmentURL != "" {
		report.AttachmentURL = input.AttachmentURL
	}
	if input.Severity != "" {
		report.Severity = input.Severity
	}

	// 5. 保存
	if err := s.repo.Update(report); err != nil {
		return nil, err
	}

	return report, nil
}

// isValidStatusTransition 校验状态流转是否合法
func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"Pending":  {"Triaged", "Closed"},
		"Triaged":  {"Resolved", "Closed"},
		"Resolved": {"Closed"},
		"Closed":   {}, // 关闭后不能再改
	}

	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
