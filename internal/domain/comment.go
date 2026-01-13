package domain

import "time"

// ReportComment 漏洞报告评论实体
type ReportComment struct {
	ID        uint      `gorm:"primaryKey;comment:评论ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 关联的漏洞报告ID
	ReportID uint `gorm:"index;not null;comment:关联的漏洞报告ID" json:"report_id"`

	// 评论作者ID
	AuthorID uint  `gorm:"index;not null;comment:评论作者ID" json:"author_id"`
	Author   *User `gorm:"-" json:"author,omitempty"` // 手动加载作者信息

	// 评论内容
	Content string `gorm:"type:text;not null;comment:评论内容" json:"content"`
}

func (ReportComment) TableName() string {
	return "report_comments"
}

// CommentRepository 评论仓库接口
type CommentRepository interface {
	Create(comment *ReportComment) error
	FindByReportID(reportID uint) ([]ReportComment, error)
	FindByID(id uint) (*ReportComment, error)
	Delete(id uint) error
}

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(reportID uint, authorID uint, content string) (*ReportComment, error)
	GetReportComments(reportID uint) ([]ReportComment, error)
	DeleteComment(commentID uint, userID uint, userRole string) error
}
