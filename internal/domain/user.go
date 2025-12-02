package domain

import (
	"time"
)

// UserService 定义了用户业务逻辑的接口
// 登录(Login) 和 注册(Register) 是业务行为，不是单纯的 CRUD
type UserService interface {
	Register(user *User) error
	Login(username, password string) (*User, string, error) // 返回用户对象和token(暂时占位)
	GetUser(id uint) (*User, error)
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
}

// TableName 指定表名和表注释
func (User) TableName() string {
	return "users"
}

// UserRepository 定义了操作数据库的接口
// 接口层只定义"我要做什么"，不关心底层是用 Postgres 还是 MySQL
type UserRepository interface {
	Create(user *User) error
	FindByUsername(username string) (*User, error)
	FindByID(id uint) (*User, error)
}
