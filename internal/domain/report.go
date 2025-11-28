package domain

import (
	"time"
)

// Report 漏洞报告实体
type Report struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 漏洞标题
	Title string `gorm:"size:255;not null" json:"title"`

	// 漏洞描述 (Text 类型，支持长文本)
	Description string `gorm:"type:text" json:""description`

	// 漏洞类型 (如: SQL Injection, XSS)
	Type string `gorm:"size:50" json:"type"`

	// 危害等级: Low, Medium, High, Critical
	Severity string `gorm:"size:20;default:'Low'" json:"severity"`

	// 状态机: Pending(待审) -> Triaged(已确) -> Resolved(已修) -> Closed(关闭)
	Status string `gorm:"size:20;default:'Pending';index" json:"status"`

	// 外键关联: 谁提交的？
	AuthorID uint `json:"author_id"`
	// GORM 会自动根据 AuthorID 去关联 User 表
	Author   User `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

// ReportRepository 接口定义
type ReportRepository interface {
	Create(report *Report) error
	FindByID(id uint) (*Report, error)
	List(page, pageSize int) ([]Report, int64, error)
	Update(report *Report) error
}

// ReportService 业务逻辑接口定义
type ReportService interface {
	SubmitReport(report *Report) error
	GetReport(id uint) (*Report, error)
	ListReports(page, pageSize int) ([]Report, int64, error)
}