package handler

import (
	"bug-bounty-lite/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
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

// UpdateProfileRequest 更新资料请求体
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

// UpdateProfile [POST] /api/v1/user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateProfile(userID.(uint), req.Name, req.Bio, req.Phone, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetProfile [GET] /api/v1/user/profile - 获取当前用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := h.Service.GetUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}

// BindOrgRequest 绑定组织请求体
type BindOrgRequest struct {
	OrgID uint `json:"org_id" binding:"required"`
}

// BindOrganization [POST] /api/v1/user/bind-org
func (h *UserHandler) BindOrganization(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req BindOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.BindOrganization(userID.(uint), req.OrgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization bound successfully"})
}

// ChangePasswordRequest 修改密码请求体
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ChangePassword [POST] /api/v1/user/change-password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}
