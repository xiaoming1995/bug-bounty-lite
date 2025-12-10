package repository

import (
	"bug-bounty-lite/internal/domain"
	"gorm.io/gorm"
)

type reportRepo struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) domain.ReportRepository {
	return &reportRepo{db: db}
}

// Create 创建报告
func (r *reportRepo) Create(report *domain.Report) error {
	// 排除关联字段，只保存基本字段和关联ID
	return r.db.Omit("Author", "Project", "VulnerabilityType", "SelfAssessment").Create(report).Error
}

// FindByID 查找详情
func (r *reportRepo) FindByID(id uint) (*domain.Report, error) {
	var report domain.Report
	// Preload 关联查询：Author、Project、VulnerabilityType、SelfAssessment
	err := r.db.Preload("Author").
		Preload("Project").
		Preload("VulnerabilityType").
		Preload("SelfAssessment").
		First(&report, id).Error
	return &report, err
}

// List 分页获取报告列表
func (r *reportRepo) List(page, pageSize int) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	// 1. 计算 Offset
	offset := (page - 1) * pageSize

	// 2. 开启一个查询会话
	query := r.db.Model(&domain.Report{})

	// 3. 先查总数 (用于前端分页条显示: 共 100 条，当前第 1 页)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 查当前页数据
	// Order("id desc") 保证最新的漏洞显示在最前面
	err := query.Preload("Author").
		Preload("Project").
		Preload("VulnerabilityType").
		Preload("SelfAssessment").
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&reports).Error

	return reports, total, err
}

// Update 更新报告状态
func (r *reportRepo) Update(report *domain.Report) error {
	// Save 会保存所有字段，排除关联字段避免错误
	return r.db.Omit("Author", "Project", "VulnerabilityType", "SelfAssessment").Save(report).Error
}