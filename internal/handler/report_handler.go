package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Service domain.ReportService
}

func NewReportHandler(s domain.ReportService) *ReportHandler {
	return &ReportHandler{Service: s}
}

// CreateReportRequest 创建报告请求 DTO
type CreateReportRequest struct {
	ProjectID           uint   `json:"project_id" binding:"required"`
	VulnerabilityName   string `json:"vulnerability_name" binding:"required,max=255"`
	VulnerabilityTypeID uint   `json:"vulnerability_type_id" binding:"required"`
	VulnerabilityImpact string `json:"vulnerability_impact"`
	SelfAssessmentID    *uint  `json:"self_assessment_id"`
	VulnerabilityURL    string `json:"vulnerability_url" binding:"omitempty,url"`
	VulnerabilityDetail string `json:"vulnerability_detail"`
	AttachmentURL       string `json:"attachment_url" binding:"omitempty,url"`
	Severity            string `json:"severity" binding:"omitempty,oneof=Low Medium High Critical"`
}

// UpdateReportRequest 更新报告请求 DTO
type UpdateReportRequest struct {
	ProjectID           uint   `json:"project_id"`
	VulnerabilityName   string `json:"vulnerability_name" binding:"omitempty,max=255"`
	VulnerabilityTypeID uint   `json:"vulnerability_type_id"`
	VulnerabilityImpact string `json:"vulnerability_impact"`
	SelfAssessmentID    *uint  `json:"self_assessment_id"`
	VulnerabilityURL    string `json:"vulnerability_url" binding:"omitempty,url"`
	VulnerabilityDetail string `json:"vulnerability_detail"`
	AttachmentURL       string `json:"attachment_url" binding:"omitempty,url"`
	Severity            string `json:"severity" binding:"omitempty,oneof=Low Medium High Critical"`
	Status              string `json:"status" binding:"omitempty,oneof=Pending Triaged Resolved Closed"`
}

// CreateHandler 提交漏洞
func (h *ReportHandler) CreateHandler(c *gin.Context) {
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 从 Context 中获取当前登录用户 ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}

	// 构建 Report 实体
	report := &domain.Report{
		ProjectID:           req.ProjectID,
		VulnerabilityName:   req.VulnerabilityName,
		VulnerabilityTypeID: req.VulnerabilityTypeID,
		VulnerabilityImpact: req.VulnerabilityImpact,
		SelfAssessmentID:    req.SelfAssessmentID,
		VulnerabilityURL:    req.VulnerabilityURL,
		VulnerabilityDetail: req.VulnerabilityDetail,
		AttachmentURL:       req.AttachmentURL,
		Severity:            req.Severity,
		AuthorID:            userID.(uint),
	}

	// 设置默认值
	if report.Severity == "" {
		report.Severity = "Low"
	}

	if err := h.Service.SubmitReport(report); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 成功时只返回提示语，不返回报告数据（统一使用 200 状态码）
	response.SuccessWithMessage(c, "漏洞报告提交成功", nil)
}

// ListHandler 获取列表
// - 白帽子只能查看自己提交的报告
// - 厂商和管理员可以查看所有报告
func (h *ReportHandler) ListHandler(c *gin.Context) {
	// 获取 query 参数 ?page=1&page_size=10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取当前用户信息
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	userID := userIDVal.(uint)

	roleVal, exists := c.Get("role")
	userRole := "whitehat" // 默认角色
	if exists && roleVal != nil {
		userRole = roleVal.(string)
	}

	keyword := c.Query("keyword")

	reports, total, err := h.Service.ListReports(page, pageSize, userID, userRole, keyword)
	if err != nil {
		response.Error(c, 500, "获取报告列表失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      reports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetHandler 获取单个详情
func (h *ReportHandler) GetHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	report, err := h.Service.GetReport(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

// UpdateHandler 更新报告
func (h *ReportHandler) UpdateHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	var req UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	// 调用 Service 更新
	report, err := h.Service.UpdateReport(uint(id), userID.(uint), role.(string), &domain.ReportUpdateInput{
		ProjectID:           req.ProjectID,
		VulnerabilityName:   req.VulnerabilityName,
		VulnerabilityTypeID: req.VulnerabilityTypeID,
		VulnerabilityImpact: req.VulnerabilityImpact,
		SelfAssessmentID:    req.SelfAssessmentID,
		VulnerabilityURL:    req.VulnerabilityURL,
		VulnerabilityDetail: req.VulnerabilityDetail,
		AttachmentURL:       req.AttachmentURL,
		Severity:            req.Severity,
		Status:              req.Status,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}

// DeleteHandler 删除报告（软删除）
// DELETE /api/v1/reports/:id
func (h *ReportHandler) DeleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的报告ID")
		return
	}

	// 获取当前用户信息
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	role, _ := c.Get("role")
	userRole := "whitehat"
	if role != nil {
		userRole = role.(string)
	}

	if err := h.Service.DeleteReport(uint(id), userID.(uint), userRole); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "报告删除成功", nil)
}

// RestoreHandler 恢复已删除的报告
// POST /api/v1/reports/:id/restore
func (h *ReportHandler) RestoreHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的报告ID")
		return
	}

	// 获取当前用户信息
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "用户未认证")
		return
	}
	role, _ := c.Get("role")
	userRole := "whitehat"
	if role != nil {
		userRole = role.(string)
	}

	if err := h.Service.RestoreReport(uint(id), userID.(uint), userRole); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "报告恢复成功", nil)
}
