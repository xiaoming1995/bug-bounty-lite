package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectTaskHandler 项目任务处理器
type ProjectTaskHandler struct {
	TaskService    domain.ProjectTaskService
	TaskRepo       domain.ProjectTaskRepository
	AssignmentRepo domain.ProjectAssignmentRepository
	ProjectRepo    domain.ProjectRepository
	AttachmentRepo domain.ProjectAttachmentRepository
}

// NewProjectTaskHandler 创建项目任务处理器
func NewProjectTaskHandler(
	taskService domain.ProjectTaskService,
	taskRepo domain.ProjectTaskRepository,
	assignmentRepo domain.ProjectAssignmentRepository,
	projectRepo domain.ProjectRepository,
	attachmentRepo domain.ProjectAttachmentRepository,
) *ProjectTaskHandler {
	return &ProjectTaskHandler{
		TaskService:    taskService,
		TaskRepo:       taskRepo,
		AssignmentRepo: assignmentRepo,
		ProjectRepo:    projectRepo,
		AttachmentRepo: attachmentRepo,
	}
}

// ListAvailableProjects 获取当前用户可见的项目列表（基于指派）
// GET /api/v1/projects/available
func (h *ProjectTaskHandler) ListAvailableProjects(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 获取用户被指派的项目
	assignments, err := h.AssignmentRepo.FindByUserID(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}

	// 定义返回结构
	type ProjectWithStatus struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Difficulty  string  `json:"difficulty"`
		Deadline    *string `json:"deadline"`
		Status      string  `json:"status"`
		Accepted    bool    `json:"accepted"`
	}

	result := make([]ProjectWithStatus, 0, len(assignments))

	for _, assignment := range assignments {
		project, err := h.ProjectRepo.FindByID(assignment.ProjectID)
		if err != nil {
			continue // 跳过不存在的项目
		}

		// 只返回招募中或进行中的项目
		if project.Status != "recruiting" && project.Status != "in_progress" {
			continue
		}

		// 检查用户是否已接受任务
		_, err = h.TaskRepo.FindByProjectAndUser(project.ID, userID.(uint))
		accepted := err == nil

		// 格式化截止日期
		var deadlineStr *string
		if project.Deadline != nil {
			formatted := project.Deadline.Format("2006-01-02")
			deadlineStr = &formatted
		}

		result = append(result, ProjectWithStatus{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Difficulty:  project.Difficulty,
			Deadline:    deadlineStr,
			Status:      project.Status,
			Accepted:    accepted,
		})
	}

	response.Success(c, gin.H{
		"list":  result,
		"total": len(result),
	})
}

// GetProjectDetail 获取项目详情（用户可见性检查）
// GET /api/v1/projects/available/:id
func (h *ProjectTaskHandler) GetProjectDetail(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 获取项目ID
	idStr := c.Param("id")
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 检查用户是否被指派到该项目
	_, err = h.AssignmentRepo.FindByProjectAndUser(uint(projectID), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusForbidden, "您无权访问该项目")
		return
	}

	// 获取项目详情
	project, err := h.ProjectRepo.FindByID(uint(projectID))
	if err != nil {
		response.Error(c, http.StatusNotFound, "项目不存在")
		return
	}

	// 检查用户是否已接受任务
	_, err = h.TaskRepo.FindByProjectAndUser(uint(projectID), userID.(uint))
	accepted := err == nil

	// 格式化截止日期
	var deadlineStr *string
	if project.Deadline != nil {
		formatted := project.Deadline.Format("2006-01-02")
		deadlineStr = &formatted
	}

	// 获取项目附件
	attachments, _ := h.AttachmentRepo.FindByProjectID(uint(projectID))

	response.Success(c, gin.H{
		"id":          project.ID,
		"name":        project.Name,
		"description": project.Description,
		"difficulty":  project.Difficulty,
		"deadline":    deadlineStr,
		"status":      project.Status,
		"created_at":  project.CreatedAt,
		"accepted":    accepted,
		"attachments": attachments,
	})
}

// AcceptTask 接受项目任务
// POST /api/v1/projects/:id/accept
func (h *ProjectTaskHandler) AcceptTask(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 获取项目ID
	idStr := c.Param("id")
	projectID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的项目ID")
		return
	}

	// 调用服务接受任务
	task, err := h.TaskService.AcceptTask(uint(projectID), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.SuccessWithMessage(c, "任务接受成功", task)
}

// ListAcceptedProjects 获取用户已接受任务的项目列表（用于漏洞提交下拉）
// GET /api/v1/projects/accepted
func (h *ProjectTaskHandler) ListAcceptedProjects(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 获取用户已接受的项目ID
	projectIDs, err := h.TaskService.GetUserAcceptedProjectIDs(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取项目列表失败")
		return
	}

	// 获取项目详情
	projects := make([]domain.Project, 0, len(projectIDs))
	for _, projectID := range projectIDs {
		project, err := h.ProjectRepo.FindByID(projectID)
		if err != nil {
			continue
		}
		projects = append(projects, *project)
	}

	response.Success(c, gin.H{
		"list":  projects,
		"total": len(projects),
	})
}
