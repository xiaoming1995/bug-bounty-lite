package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SystemConfigHandler struct {
	Service domain.SystemConfigService
}

func NewSystemConfigHandler(s domain.SystemConfigService) *SystemConfigHandler {
	return &SystemConfigHandler{Service: s}
}

// CreateConfigRequest 创建配置请求 DTO
type CreateConfigRequest struct {
	ConfigType  string          `json:"config_type" binding:"required,max=50"`
	ConfigKey   string          `json:"config_key" binding:"required,max=100"`
	ConfigValue string          `json:"config_value" binding:"required,max=255"`
	Description string          `json:"description"`
	SortOrder   int             `json:"sort_order"`
	Status      string          `json:"status" binding:"omitempty,oneof=active inactive"`
	ExtraData   domain.JSON     `json:"extra_data"`
}

// UpdateConfigRequest 更新配置请求 DTO
type UpdateConfigRequest struct {
	ConfigType  string          `json:"config_type" binding:"omitempty,max=50"`
	ConfigKey   string          `json:"config_key" binding:"omitempty,max=100"`
	ConfigValue string          `json:"config_value" binding:"omitempty,max=255"`
	Description string          `json:"description"`
	SortOrder   int             `json:"sort_order"`
	Status      string          `json:"status" binding:"omitempty,oneof=active inactive"`
	ExtraData   domain.JSON     `json:"extra_data"`
}

// GetConfigsByTypeHandler 根据类型获取配置列表
// GET /api/v1/configs/:type
func (h *SystemConfigHandler) GetConfigsByTypeHandler(c *gin.Context) {
	configType := c.Param("type")
	if configType == "" {
		response.Error(c, http.StatusBadRequest, "配置类型不能为空")
		return
	}

	// 获取当前用户角色，判断是否包含非活跃配置
	role, exists := c.Get("role")
	includeInactive := exists && role.(string) == "admin"

	configs, err := h.Service.GetConfigsByType(configType, includeInactive)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, configs)
}

// GetConfigHandler 获取配置详情
// GET /api/v1/configs/:type/:id
func (h *SystemConfigHandler) GetConfigHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置ID")
		return
	}

	config, err := h.Service.GetConfig(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, config)
}

// CreateConfigHandler 创建配置
// POST /api/v1/configs/:type
func (h *SystemConfigHandler) CreateConfigHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以创建配置
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以创建配置")
		return
	}

	configType := c.Param("type")
	if configType == "" {
		response.Error(c, http.StatusBadRequest, "配置类型不能为空")
		return
	}

	var req CreateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 确保 config_type 与路径参数一致
	req.ConfigType = configType

	config := &domain.SystemConfig{
		ConfigType:  req.ConfigType,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		ExtraData:   req.ExtraData,
	}

	if config.Status == "" {
		config.Status = "active"
	}

	if err := h.Service.CreateConfig(config); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, config)
}

// UpdateConfigHandler 更新配置
// PUT /api/v1/configs/:type/:id
func (h *SystemConfigHandler) UpdateConfigHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以更新配置
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以更新配置")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置ID")
		return
	}

	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	config := &domain.SystemConfig{
		ConfigType:  req.ConfigType,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		ExtraData:   req.ExtraData,
	}

	if err := h.Service.UpdateConfig(uint(id), config); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedConfig, _ := h.Service.GetConfig(uint(id))
	response.Success(c, updatedConfig)
}

// DeleteConfigHandler 删除配置
// DELETE /api/v1/configs/:type/:id
func (h *SystemConfigHandler) DeleteConfigHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以删除配置
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以删除配置")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的配置ID")
		return
	}

	if err := h.Service.DeleteConfig(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "配置删除成功", nil)
}

