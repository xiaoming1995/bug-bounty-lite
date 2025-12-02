package domain

import (
	"time"
)

// Report 漏洞报告实体
type Report struct {
	ID        uint      `gorm:"primaryKey;comment:报告ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 漏洞标题
	Title string `gorm:"size:255;not null;comment:漏洞标题" json:"title"`

	// 漏洞描述 (Text 类型，支持长文本)
	Description string `gorm:"type:text;comment:漏洞描述" json:"description"`

	// 漏洞类型 (如: SQL Injection, XSS)
	Type string `gorm:"size:50;comment:漏洞类型(如XSS/SQLi/CSRF)" json:"type"`

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
	Title       string
	Description string
	Type        string
	Severity    string
	Status      string
}

// ReportService 业务逻辑接口定义
type ReportService interface {
	SubmitReport(report *Report) error
	GetReport(id uint) (*Report, error)
	ListReports(page, pageSize int) ([]Report, int64, error)
	UpdateReport(id uint, userID uint, userRole string, input *ReportUpdateInput) (*Report, error)
}
