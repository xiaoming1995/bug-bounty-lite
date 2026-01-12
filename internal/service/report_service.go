package service

import (
	"bug-bounty-lite/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type reportService struct {
	repo             domain.ReportRepository
	systemConfigRepo domain.SystemConfigRepository
}

func NewReportService(repo domain.ReportRepository, systemConfigRepo domain.SystemConfigRepository) domain.ReportService {
	return &reportService{
		repo:             repo,
		systemConfigRepo: systemConfigRepo,
	}
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

	// 3. 验证危害自评ID（如果提供了）
	if report.SelfAssessmentID != nil && *report.SelfAssessmentID != 0 {
		config, err := s.systemConfigRepo.FindByID(*report.SelfAssessmentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("危害自评配置ID不存在")
			}
			return errors.New("验证危害自评配置时出错: " + err.Error())
		}
		if config.ConfigType != "severity_level" {
			return errors.New("危害自评配置ID必须是危害等级类型")
		}
	}

	// 4. 设置默认值
	if report.Severity == "" {
		report.Severity = "Low"
	}

	// 5. 调用 Repo 创建
	return s.repo.Create(report)
}

func (s *reportService) GetReport(id uint) (*domain.Report, error) {
	return s.repo.FindByID(id)
}

// ListReports 获取报告列表
// - 白帽子只能查看自己提交的报告
// - 厂商和管理员可以查看所有报告
func (s *reportService) ListReports(page, pageSize int, userID uint, userRole string, keyword string) ([]domain.Report, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 根据角色决定查询范围
	var authorID *uint
	if userRole == "whitehat" {
		// 白帽子只能查看自己的报告
		authorID = &userID
	}
	// 厂商和管理员可以查看所有报告，authorID 为 nil

	return s.repo.List(page, pageSize, authorID, keyword)
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
	if input.SelfAssessmentID != nil {
		if *input.SelfAssessmentID != 0 {
			// 验证危害自评ID
			config, err := s.systemConfigRepo.FindByID(*input.SelfAssessmentID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("危害自评配置ID不存在")
				}
				return nil, errors.New("验证危害自评配置时出错: " + err.Error())
			}
			if config.ConfigType != "severity_level" {
				return nil, errors.New("危害自评配置ID必须是危害等级类型")
			}
			report.SelfAssessmentID = input.SelfAssessmentID
		} else {
			// 如果明确传了 0，设置为 nil（表示数据库中的 NULL）
			report.SelfAssessmentID = nil
		}
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

// DeleteReport 软删除报告
func (s *reportService) DeleteReport(id uint, userID uint, userRole string) error {
	// 1. 获取报告
	report, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("报告不存在")
	}

	// 2. 权限校验：只有报告作者或管理员可以删除
	if report.AuthorID != userID && userRole != "admin" {
		return errors.New("没有权限删除此报告")
	}

	// 3. 执行软删除
	return s.repo.Delete(id)
}

// RestoreReport 恢复已删除的报告
func (s *reportService) RestoreReport(id uint, userID uint, userRole string) error {
	// 1. 只有管理员可以恢复报告
	if userRole != "admin" {
		return errors.New("只有管理员可以恢复报告")
	}

	// 2. 检查报告是否存在（包含已删除的）
	report, err := s.repo.FindByIDWithDeleted(id)
	if err != nil {
		return errors.New("报告不存在")
	}

	// 3. 检查报告是否已被删除
	if !report.DeletedAt.Valid {
		return errors.New("报告未被删除，无需恢复")
	}

	// 4. 执行恢复
	return s.repo.Restore(id)
}
