package service

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/jwt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo       domain.UserRepository
	jwtManager *jwt.JWTManager
}

// NewUserService 构造函数
func NewUserService(repo domain.UserRepository, jwtManager *jwt.JWTManager) domain.UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register 用户注册
// 核心逻辑：接收明文密码 -> bcrypt加密 -> 存入数据库
func (s *userService) Register(user *domain.User) error {
	// 1. 检查用户名是否已存在
	existingUser, _ := s.repo.FindByUsername(user.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// 2. 密码加密 (Cost 默认 10)
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. 用加密后的哈希值替换原始密码
	user.Password = string(hashedPwd)

	// 4. 调用 Repo 保存到数据库
	return s.repo.Create(user)
}

// Login 用户登录
// 核心逻辑：接收明文密码 -> bcrypt比对数据库里的哈希 -> 成功则返回用户和JWT
func (s *userService) Login(username, password string) (*domain.User, string, error) {
	// 1. 根据用户名找用户
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		// 返回模糊错误，防止攻击者枚举用户名
		return nil, "", errors.New("invalid username or password")
	}

	// 2. 比对密码
	// 注意：这里必须用 CompareHashAndPassword，不能用 ==
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 3. 生成 JWT Token
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

// GetUser 获取用户信息
func (s *userService) GetUser(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}
