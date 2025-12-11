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
// authorID 为 nil 时查询所有报告，否则只查询指定用户的报告
func (r *reportRepo) List(page, pageSize int, authorID *uint) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	// 1. 计算 Offset
	offset := (page - 1) * pageSize

	// 2. 查询总数
	countQuery := r.db.Model(&domain.Report{})
	if authorID != nil && *authorID > 0 {
		countQuery = countQuery.Where("author_id = ?", *authorID)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. 查询当前页数据
	dataQuery := r.db.Preload("Author").
		Preload("Project").
		Preload("VulnerabilityType").
		Preload("SelfAssessment")

	// 添加 author_id 过滤条件
	if authorID != nil && *authorID > 0 {
		dataQuery = dataQuery.Where("author_id = ?", *authorID)
	}

	// Order("id desc") 保证最新的漏洞显示在最前面
	err := dataQuery.Order("id desc").
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