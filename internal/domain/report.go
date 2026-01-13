package domain

import (
	"bug-bounty-lite/pkg/types"

	"gorm.io/gorm"
)

// Report 漏洞报告实体
type Report struct {
	ID        uint           `gorm:"primaryKey;comment:报告ID" json:"id"`
	CreatedAt types.DateTime `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt types.DateTime `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`

	// 关联项目（必填）- 不使用数据库外键，使用代码逻辑验证
	ProjectID uint    `gorm:"not null;index;comment:关联项目ID(必填)" json:"project_id"`
	Project   Project `gorm:"-" json:"project,omitempty"` // 不创建外键，手动加载

	// 漏洞名称（必填）
	VulnerabilityName string `gorm:"size:255;not null;comment:漏洞名称(必填，文本输入)" json:"vulnerability_name"`

	// 关联漏洞类型配置（必填）- 不使用数据库外键
	VulnerabilityTypeID uint         `gorm:"not null;index;comment:关联漏洞类型配置ID(必填)" json:"vulnerability_type_id"`
	VulnerabilityType   SystemConfig `gorm:"-" json:"vulnerability_type,omitempty"` // 不创建外键，手动加载

	// 漏洞的危害
	VulnerabilityImpact string `gorm:"type:text;comment:漏洞的危害(文本输入，描述漏洞可能造成的危害)" json:"vulnerability_impact"`

	// 危害自评（关联危害等级配置）- 不使用数据库外键
	SelfAssessmentID *uint        `gorm:"column:self_assessment_id;index;comment:危害自评ID(关联config表)" json:"self_assessment_id"`
	SelfAssessment   SystemConfig `gorm:"-" json:"self_assessment,omitempty"` // 不创建外键，手动加载

	// 漏洞链接
	VulnerabilityURL string `gorm:"size:500;comment:漏洞链接(URL格式，指向漏洞相关页面)" json:"vulnerability_url"`

	// 漏洞详情
	VulnerabilityDetail string `gorm:"type:text;comment:漏洞详情(文本输入，详细描述漏洞情况)" json:"vulnerability_detail"`

	// 附件地址
	AttachmentURL string `gorm:"size:500;comment:附件地址(文件上传后的URL，单个文件，后续可扩展为多个)" json:"attachment_url"`

	// 危害等级: Low, Medium, High, Critical
	Severity string `gorm:"size:20;comment:危害等级(Critical:严重, High:高危, Medium:中危, Low:低危, None:无危害)" json:"severity"`

	// 状态机: Pending(待审核) -> Audited(已审核) / Rejected(驳回)
	Status string `gorm:"size:20;default:'Pending';index;comment:报告状态(Pending:待审核[默认], Audited:已审核, Rejected:驳回)" json:"status"`

	// 提交者ID - 不使用数据库外键
	AuthorID uint `gorm:"index;comment:提交者ID" json:"author_id"`
	Author   User `gorm:"-" json:"author,omitempty"` // 不创建外键，手动加载
}

// TableName 指定表名
func (Report) TableName() string {
	return "reports"
}

// ReportRepository 接口定义
type ReportRepository interface {
	Create(report *Report) error
	FindByID(id uint) (*Report, error)
	FindByIDWithDeleted(id uint) (*Report, error) // 包含已删除的报告
	List(page, pageSize int, authorID *uint, keyword string) ([]Report, int64, error)
	Update(report *Report) error
	Delete(id uint) error  // 软删除
	Restore(id uint) error // 恢复已删除的报告
}

// ReportUpdateInput 更新报告输入
type ReportUpdateInput struct {
	ProjectID           uint
	VulnerabilityName   string
	VulnerabilityTypeID uint
	VulnerabilityImpact string
	SelfAssessmentID    *uint
	VulnerabilityURL    string
	VulnerabilityDetail string
	AttachmentURL       string
	Severity            string
	Status              string
}

// ReportService 业务逻辑接口定义
type ReportService interface {
	SubmitReport(report *Report) error
	GetReport(id uint) (*Report, error)
	ListReports(page, pageSize int, userID uint, userRole string, keyword string) ([]Report, int64, error)
	UpdateReport(id uint, userID uint, userRole string, input *ReportUpdateInput) (*Report, error)
	DeleteReport(id uint, userID uint, userRole string) error  // 软删除
	RestoreReport(id uint, userID uint, userRole string) error // 恢复已删除的报告
}
