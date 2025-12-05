package domain

import (
	"time"
)

// Report 漏洞报告实体
type Report struct {
	ID        uint      `gorm:"primaryKey;comment:报告ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 关联项目（必填）
	ProjectID uint `gorm:"not null;index;comment:关联项目ID(必填，关联projects表)" json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`

	// 漏洞名称（必填）
	VulnerabilityName string `gorm:"size:255;not null;comment:漏洞名称(必填，文本输入)" json:"vulnerability_name"`

	// 关联漏洞类型配置（必填）
	VulnerabilityTypeID uint `gorm:"not null;index;comment:关联漏洞类型配置ID(必填，关联system_configs表，config_type='vulnerability_type')" json:"vulnerability_type_id"`
	VulnerabilityType   SystemConfig `gorm:"foreignKey:VulnerabilityTypeID" json:"vulnerability_type,omitempty"`

	// 漏洞的危害
	VulnerabilityImpact string `gorm:"type:text;comment:漏洞的危害(文本输入，描述漏洞可能造成的危害)" json:"vulnerability_impact"`

	// 危害自评
	SelfAssessment string `gorm:"type:text;comment:危害自评(文本输入，提交者对漏洞危害的自我评估)" json:"self_assessment"`

	// 漏洞链接
	VulnerabilityURL string `gorm:"size:500;comment:漏洞链接(URL格式，指向漏洞相关页面)" json:"vulnerability_url"`

	// 漏洞详情
	VulnerabilityDetail string `gorm:"type:text;comment:漏洞详情(文本输入，详细描述漏洞情况)" json:"vulnerability_detail"`

	// 附件地址
	AttachmentURL string `gorm:"size:500;comment:附件地址(文件上传后的URL，单个文件，后续可扩展为多个)" json:"attachment_url"`

	// ========== 保留字段（向后兼容） ==========
	// 漏洞标题（保留字段，与vulnerability_name同步）
	Title string `gorm:"size:255;not null;comment:漏洞标题(保留字段，与vulnerability_name同步，用于向后兼容)" json:"title"`

	// 漏洞描述（保留字段，与vulnerability_detail同步）
	Description string `gorm:"type:text;comment:漏洞描述(保留字段，与vulnerability_detail同步，用于向后兼容)" json:"description"`

	// 漏洞类型（保留字段，从vulnerability_type配置同步）
	Type string `gorm:"size:50;comment:漏洞类型(保留字段，从vulnerability_type配置同步，用于向后兼容)" json:"type"`

	// 危害等级: Low, Medium, High, Critical
	Severity string `gorm:"size:20;default:'Low';comment:危害等级(Low/Medium/High/Critical)" json:"severity"`

	// 状态机: Pending(待审) -> Triaged(已确) -> Resolved(已修) -> Closed(关闭)
	Status string `gorm:"size:20;default:'Pending';index;comment:报告状态(Pending/Triaged/Resolved/Closed)" json:"status"`

	// 外键关联: 谁提交的？
	AuthorID uint `gorm:"comment:提交者ID" json:"author_id"`
	// GORM 会自动根据 AuthorID 去关联 User 表
	Author User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// TableName 指定表名
func (Report) TableName() string {
	return "reports"
}

// ReportRepository 接口定义
type ReportRepository interface {
	Create(report *Report) error
	FindByID(id uint) (*Report, error)
	List(page, pageSize int) ([]Report, int64, error)
	Update(report *Report) error
}

// ReportUpdateInput 更新报告输入
type ReportUpdateInput struct {
	ProjectID           uint
	VulnerabilityName    string
	VulnerabilityTypeID  uint
	VulnerabilityImpact string
	SelfAssessment      string
	VulnerabilityURL    string
	VulnerabilityDetail string
	AttachmentURL       string
	Title               string
	Description         string
	Type                string
	Severity            string
	Status              string
}

// ReportService 业务逻辑接口定义
type ReportService interface {
	SubmitReport(report *Report) error
	GetReport(id uint) (*Report, error)
	ListReports(page, pageSize int) ([]Report, int64, error)
	UpdateReport(id uint, userID uint, userRole string, input *ReportUpdateInput) (*Report, error)
}
