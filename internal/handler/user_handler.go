package handler

import (
	"bug-bounty-lite/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	Service domain.UserService
}

func NewUserHandler(s domain.UserService) *UserHandler {
	return &UserHandler{Service: s}
}

// DTO: 专门用于接收注册参数
type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"` // 密码最少6位
}

// Register 处理注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	// 1. 解析并校验 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 转换为 Domain 实体
	user := &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	// 3. 调用 Service
	if err := h.Service.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// DTO: 专门用于接收登录参数
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 处理登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Service 进行登录
	user, token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		// 登录失败返回 401 Unauthorized
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user, // 注意：User 里的 Password 字段有 `json:"-"`，所以不会返回
	})
}