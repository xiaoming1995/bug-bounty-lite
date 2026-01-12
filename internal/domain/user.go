package domain

import (
	"time"
)

// UserService 定义了用户业务逻辑的接口
// 登录(Login) 和 注册(Register) 是业务行为，不是单纯的 CRUD
type UserService interface {
	Register(user *User) error
	Login(username, password string) (*User, string, error)
	GetUser(id uint) (*User, error)
	UpdateProfile(userID uint, name string, bio string, phone string, email string) error // 更新基本信息与简介
	ChangePassword(userID uint, oldPassword, newPassword string) error                    // 修改密码
	BindOrganization(userID uint, orgID uint) error                                       // 绑定组织
}

// OrganizationService 组织业务接口
type OrganizationService interface {
	CreateOrganization(name string, description string) (*Organization, error)
	GetOrganization(id uint) (*Organization, error)
	ListOrganizations() ([]Organization, error)
	UpdateOrganization(id uint, name string, description string) (*Organization, error)
	DeleteOrganization(id uint) error
}

// User 用户实体
type User struct {
	ID        uint      `gorm:"primaryKey;comment:用户ID" json:"id"`
	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`

	// uniqueIndex 保证用户名唯一
	Username string `gorm:"size:64;uniqueIndex;not null;comment:用户名" json:"username"`

	// json:"-" 表示该字段永远不会被 JSON 序列化返回给前端（密码绝不能泄露）
	Password string `gorm:"size:255;not null;comment:密码(bcrypt加密)" json:"-"`

	// 角色: "whitehat" (白帽子), "vendor" (厂商), "admin" (管理员)
	Role string `gorm:"size:20;default:'whitehat';comment:用户角色(whitehat/vendor/admin)" json:"role"`

	// 用户信息（需要审核后才能生效）
	Phone string `gorm:"size:20;comment:手机号" json:"phone"`
	Email string `gorm:"size:100;comment:邮箱" json:"email"`
	Name  string `gorm:"size:50;comment:姓名" json:"name"`

	// --- 新增字段 ---
	Bio   string        `gorm:"type:text;comment:个人简介" json:"bio"`
	OrgID uint          `gorm:"index;comment:所属组织ID" json:"org_id"`
	Org   *Organization `gorm:"-" json:"org,omitempty"` // 移除自动外键，改用代码逻辑手动加载

	LastLoginAt *time.Time `gorm:"comment:最后登录时间" json:"last_login_at"`
}

// Organization 组织实体
type Organization struct {
	ID          uint      `gorm:"primaryKey;comment:组织ID" json:"id"`
	CreatedAt   time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"comment:更新时间" json:"updated_at"`
	Name        string    `gorm:"size:100;uniqueIndex;not null;comment:组织名称" json:"name"`
	Description string    `gorm:"type:text;comment:组织描述" json:"description"`
}

func (Organization) TableName() string {
	return "organizations"
}

// TableName 指定表名和表注释
func (User) TableName() string {
	return "users"
}

// UserRepository 定义了操作数据库的接口
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	UpdateLastLoginAt(userID uint, loginTime time.Time) error
	UpdateProfileFields(userID uint, name, bio, phone, email string) error
	FindByUsername(username string) (*User, error)
	FindByID(id uint) (*User, error)
}

// OrganizationRepository 组织仓库接口
type OrganizationRepository interface {
	Create(org *Organization) error
	FindByID(id uint) (*Organization, error)
	List() ([]Organization, error)
	Update(org *Organization) error
	Delete(id uint) error
}
