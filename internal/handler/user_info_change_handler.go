package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserInfoChangeHandler struct {
	Service domain.UserInfoChangeService
}

func NewUserInfoChangeHandler(s domain.UserInfoChangeService) *UserInfoChangeHandler {
	return &UserInfoChangeHandler{Service: s}
}

// SubmitChangeRequestRequest 提交变更申请的请求体
type SubmitChangeRequestRequest struct {
	Phone string `json:"phone" binding:"omitempty"` // 可选
	Email string `json:"email" binding:"omitempty,email"` // 可选，如果提供则必须是邮箱格式
	Name  string `json:"name" binding:"omitempty"`   // 可选
}

// SubmitChangeRequest 提交用户信息变更申请
// POST /api/v1/user/info/change
func (h *UserInfoChangeHandler) SubmitChangeRequest(c *gin.Context) {
	// 1. 从上下文获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 2. 解析请求体
	var req SubmitChangeRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 3. 验证至少提供一个字段
	if req.Phone == "" && req.Email == "" && req.Name == "" {
		response.Error(c, http.StatusBadRequest, "至少需要提供一个要变更的字段（手机号、邮箱或姓名）")
		return
	}

	// 4. 调用服务层
	request, err := h.Service.SubmitChangeRequest(
		userID.(uint),
		req.Phone,
		req.Email,
		req.Name,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 5. 返回成功响应
	response.Created(c, request)
}

// GetUserChangeRequests 获取用户的所有变更申请
// GET /api/v1/user/info/changes
func (h *UserInfoChangeHandler) GetUserChangeRequests(c *gin.Context) {
	// 1. 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 2. 调用服务层
	requests, err := h.Service.GetUserChangeRequests(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取变更申请列表失败")
		return
	}

	// 3. 返回成功响应
	response.SuccessWithMessage(c, "获取成功", requests)
}

// GetChangeRequest 获取单个变更申请详情
// GET /api/v1/user/info/changes/:id
func (h *UserInfoChangeHandler) GetChangeRequest(c *gin.Context) {
	// 1. 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 2. 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的申请ID")
		return
	}

	// 3. 调用服务层
	request, err := h.Service.GetChangeRequest(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. 返回成功响应
	response.SuccessWithMessage(c, "获取成功", request)
}

