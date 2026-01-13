package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"bug-bounty-lite/pkg/upload"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AvatarHandler 头像管理处理器
type AvatarHandler struct {
	service domain.AvatarService
}

// NewAvatarHandler 创建头像处理器实例
func NewAvatarHandler(s domain.AvatarService) *AvatarHandler {
	return &AvatarHandler{service: s}
}

// UploadAvatarHandler 管理员上传头像到平台头像库
// POST /api/v1/avatars/upload
func (h *AvatarHandler) UploadAvatarHandler(c *gin.Context) {
	// 检查是否为管理员
	role, _ := c.Get("role")
	if role != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以上传头像")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的头像文件")
		return
	}

	// 获取头像名称
	name := c.PostForm("name")
	if name == "" {
		name = file.Filename
	}

	// 构建基础 URL
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	// 上传文件到 avatars 子目录
	result, err := upload.UploadFileToDir(file, baseURL, "avatars")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 保存头像记录到数据库
	avatar, err := h.service.UploadAvatar(name, result.URL)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "保存头像信息失败")
		return
	}

	response.Success(c, avatar)
}

// ListAvatarsHandler 获取所有头像（管理员）
// GET /api/v1/avatars
func (h *AvatarHandler) ListAvatarsHandler(c *gin.Context) {
	avatars, err := h.service.ListAvatars()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取头像列表失败")
		return
	}
	response.Success(c, avatars)
}

// ListActiveAvatarsHandler 获取启用的头像（用户选择用）
// GET /api/v1/avatars/active
func (h *AvatarHandler) ListActiveAvatarsHandler(c *gin.Context) {
	avatars, err := h.service.ListActiveAvatars()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取头像列表失败")
		return
	}
	response.Success(c, avatars)
}

// AdminUpdateAvatarRequest 管理员更新头像请求体
type AdminUpdateAvatarRequest struct {
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	SortOrder int    `json:"sort_order"`
}

// UpdateAvatarHandler 更新头像信息（管理员）
// PUT /api/v1/avatars/:id
func (h *AvatarHandler) UpdateAvatarHandler(c *gin.Context) {
	// 检查是否为管理员
	role, _ := c.Get("role")
	if role != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以修改头像信息")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的头像ID")
		return
	}

	var req AdminUpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	avatar, err := h.service.UpdateAvatar(uint(id), req.Name, req.IsActive, req.SortOrder)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "更新头像失败")
		return
	}
	response.Success(c, avatar)
}

// DeleteAvatarHandler 删除头像（管理员）
// DELETE /api/v1/avatars/:id
func (h *AvatarHandler) DeleteAvatarHandler(c *gin.Context) {
	// 检查是否为管理员
	role, _ := c.Get("role")
	if role != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以删除头像")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的头像ID")
		return
	}

	if err := h.service.DeleteAvatar(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除头像失败")
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
