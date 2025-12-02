package domain

import (
	"time"
)

// UserInfoChangeRequest 用户信息变更申请表
type UserInfoChangeRequest struct {
	ID        uint      `gorm:"primaryKey;comment:申请ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// 关联用户
	UserID uint `gorm:"not null;index;comment:用户ID" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 待审核的信息
	Phone string `gorm:"size:20;comment:手机号" json:"phone"`
	Email string `gorm:"size:100;comment:邮箱" json:"email"`
	Name  string `gorm:"size:50;comment:姓名" json:"name"`

	// 审核状态: pending(待审核), approved(已通过), rejected(已拒绝)
	Status string `gorm:"size:20;default:'pending';index;comment:审核状态(pending/approved/rejected)" json:"status"`

	// 审核信息（后台审核时填写）
	ReviewedAt  *time.Time `gorm:"comment:审核时间" json:"reviewed_at,omitempty"`
	ReviewerID  *uint      `gorm:"comment:审核人ID" json:"reviewer_id,omitempty"`
	ReviewNote  string     `gorm:"type:text;comment:审核备注" json:"review_note,omitempty"`
}

// TableName 指定表名
func (UserInfoChangeRequest) TableName() string {
	return "user_info_change_requests"
}

// UserInfoChangeRepository 用户信息变更申请仓库接口
type UserInfoChangeRepository interface {
	Create(request *UserInfoChangeRequest) error
	FindByID(id uint) (*UserInfoChangeRequest, error)
	FindByUserID(userID uint) ([]UserInfoChangeRequest, error)
	FindPendingByUserID(userID uint) (*UserInfoChangeRequest, error)
	Update(request *UserInfoChangeRequest) error
}

// UserInfoChangeService 用户信息变更服务接口
type UserInfoChangeService interface {
	SubmitChangeRequest(userID uint, phone, email, name string) (*UserInfoChangeRequest, error)
	GetUserChangeRequests(userID uint) ([]UserInfoChangeRequest, error)
	GetChangeRequest(id uint, userID uint) (*UserInfoChangeRequest, error)
}

