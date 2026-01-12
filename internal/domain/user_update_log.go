package domain

import "time"

// UserUpdateLog 用户信息修改记录
type UserUpdateLog struct {
	ID        uint      `gorm:"primaryKey;comment:记录ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:记录时间" json:"created_at"`

	UserID uint   `gorm:"not null;index;comment:用户ID" json:"user_id"`
	Field  string `gorm:"size:50;not null;comment:修改字段" json:"field"`
	Before string `gorm:"type:text;comment:修改前的值" json:"before"`
	After  string `gorm:"type:text;comment:修改后的值" json:"after"`
	Reason string `gorm:"size:255;comment:修改原因" json:"reason"`
}

func (UserUpdateLog) TableName() string {
	return "user_update_logs"
}

// UserUpdateLogRepository 修改记录仓库接口
type UserUpdateLogRepository interface {
	Create(log *UserUpdateLog) error
	FindByUserID(userID uint) ([]UserUpdateLog, error)
	GetLastUpdateAt(userID uint, field string) (*time.Time, error)
}
