package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	Service domain.ProjectService
}

func NewProjectHandler(s domain.ProjectService) *ProjectHandler {
	return &ProjectHandler{Service: s}
}

// CreateProjectRequest 创建项目请求 DTO
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description"`
	Note        string `json:"note"`
}

// UpdateProjectRequest 更新项目请求 DTO
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"omitempty,max=255"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// CreateHandler 创建项目
// POST /api/v1/projects
func (h *ProjectHandler) CreateHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以创建项目
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以创建项目")
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 构建 Project 实体
	project := &domain.Project{
		Name:        req.Name,
		Description: req.Description,
		Note:        req.Note,
		Status:      "active", // 默认状态
	}

	if err := h.Service.CreateProject(project); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, project)
}

// ListHandler 获取项目列表
// GET /api/v1/projects
func (h *ProjectHandler) ListHandler(c *gin.Context) {
	// 获取 query 参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取当前用户角色，判断是否包含非活跃项目
	role, exists := c.Get("role")
	includeInactive := exists && role.(string) == "admin"

	projects, total, err := h.Service.ListProjects(page, pageSize, includeInactive)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      projects,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetHandler 获取项目详情
// GET /api/v1/projects/:id
func (h *ProjectHandler) GetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 获取当前用户角色，判断是否包含非活跃项目
	role, exists := c.Get("role")
	includeInactive := exists && role.(string) == "admin"

	project, err := h.Service.GetProject(uint(id), includeInactive)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, project)
}

// UpdateHandler 更新项目
// PUT /api/v1/projects/:id
func (h *ProjectHandler) UpdateHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以更新项目
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以更新项目")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 调用 Service 更新
	project, err := h.Service.UpdateProject(uint(id), &domain.ProjectUpdateInput{
		Name:        req.Name,
		Description: req.Description,
		Note:        req.Note,
		Status:      req.Status,
	})

	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, project)
}

// DeleteHandler 删除项目
// DELETE /api/v1/projects/:id
func (h *ProjectHandler) DeleteHandler(c *gin.Context) {
	// 权限检查：只有 admin 可以删除项目
	role, exists := c.Get("role")
	if !exists || role.(string) != "admin" {
		response.Error(c, http.StatusForbidden, "只有管理员可以删除项目")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	if err := h.Service.DeleteProject(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "项目删除成功", nil)
}

