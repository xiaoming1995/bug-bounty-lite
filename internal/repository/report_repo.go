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
	return r.db.Create(report).Error
}

// loadAssociations 手动加载关联数据（因为移除了外键）
func (r *reportRepo) loadAssociations(report *domain.Report) error {
	// 加载 Author
	if report.AuthorID > 0 {
		var author domain.User
		if err := r.db.First(&author, report.AuthorID).Error; err == nil {
			report.Author = author
		}
	}

	// 加载 Project（使用 Unscoped 以防项目被软删除）
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

	return nil
}

// FindByID 查找详情（不包含已删除的）
func (r *reportRepo) FindByID(id uint) (*domain.Report, error) {
	var report domain.Report
	if err := r.db.First(&report, id).Error; err != nil {
		return nil, err
	}
	// 手动加载关联数据
	r.loadAssociations(&report)
	return &report, nil
}

// FindByIDWithDeleted 查找详情（包含已删除的）
func (r *reportRepo) FindByIDWithDeleted(id uint) (*domain.Report, error) {
	var report domain.Report
	if err := r.db.Unscoped().First(&report, id).Error; err != nil {
		return nil, err
	}
	// 手动加载关联数据
	r.loadAssociations(&report)
	return &report, nil
}

// List 分页获取报告列表
// authorID 为 nil 时查询所有报告，否则只查询指定用户的报告
func (r *reportRepo) List(page, pageSize int, authorID *uint, keyword string) ([]domain.Report, int64, error) {
	var reports []domain.Report
	var total int64

	// 1. 计算 Offset
	offset := (page - 1) * pageSize

	// 2. 查询总数
	baseQuery := r.db.Model(&domain.Report{})
	if authorID != nil && *authorID > 0 {
		baseQuery = baseQuery.Where("author_id = ?", *authorID)
	}
	if keyword != "" {
		baseQuery = baseQuery.Where("vulnerability_name LIKE ?", "%"+keyword+"%")
	}

	if err := baseQuery.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. 查询当前页数据
	if err := baseQuery.Session(&gorm.Session{}).Order("id desc").Offset(offset).Limit(pageSize).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	// 4. 手动加载每个报告的关联数据
	for i := range reports {
		r.loadAssociations(&reports[i])
	}

	return reports, total, nil
}

// Update 更新报告
func (r *reportRepo) Update(report *domain.Report) error {
	return r.db.Save(report).Error
}

// Delete 软删除报告
func (r *reportRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Report{}, id).Error
}

// Restore 恢复已删除的报告
func (r *reportRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&domain.Report{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
