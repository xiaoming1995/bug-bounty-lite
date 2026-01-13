package domain

import "time"

// Avatar 头像实体
type Avatar struct {
	ID        uint      `gorm:"primaryKey;comment:头像ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	Name      string `gorm:"size:100;comment:头像名称" json:"name"`
	URL       string `gorm:"size:500;not null;comment:头像URL" json:"url"`
	IsActive  bool   `gorm:"default:true;comment:是否启用" json:"is_active"`
	SortOrder int    `gorm:"default:0;comment:排序" json:"sort_order"`
}

func (Avatar) TableName() string {
	return "avatars"
}

// AvatarRepository 头像仓库接口
type AvatarRepository interface {
	Create(avatar *Avatar) error
	FindByID(id uint) (*Avatar, error)
	List() ([]Avatar, error)
	ListActive() ([]Avatar, error)
	Update(avatar *Avatar) error
	Delete(id uint) error
}

// AvatarService 头像服务接口
type AvatarService interface {
	UploadAvatar(name string, url string) (*Avatar, error)
	GetAvatar(id uint) (*Avatar, error)
	ListAvatars() ([]Avatar, error)
	ListActiveAvatars() ([]Avatar, error)
	UpdateAvatar(id uint, name string, isActive bool, sortOrder int) (*Avatar, error)
	DeleteAvatar(id uint) error
}
