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
	SelfAssessment      string `json:"self_assessment"`
	VulnerabilityURL    string `json:"vulnerability_url" binding:"omitempty,url"`
	VulnerabilityDetail string `json:"vulnerability_detail"`
	AttachmentURL       string `json:"attachment_url" binding:"omitempty,url"`
	Severity            string `json:"severity" binding:"omitempty,oneof=Low Medium High Critical"`
	// 保留字段（向后兼容）
	Title       string `json:"title" binding:"omitempty,max=255"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"omitempty,max=50"`
}

// UpdateReportRequest 更新报告请求 DTO
type UpdateReportRequest struct {
	ProjectID           uint   `json:"project_id"`
	VulnerabilityName   string `json:"vulnerability_name" binding:"omitempty,max=255"`
	VulnerabilityTypeID uint   `json:"vulnerability_type_id"`
	VulnerabilityImpact string `json:"vulnerability_impact"`
	SelfAssessment      string `json:"self_assessment"`
	VulnerabilityURL    string `json:"vulnerability_url" binding:"omitempty,url"`
	VulnerabilityDetail string `json:"vulnerability_detail"`
	AttachmentURL       string `json:"attachment_url" binding:"omitempty,url"`
	Severity            string `json:"severity" binding:"omitempty,oneof=Low Medium High Critical"`
	Status              string `json:"status" binding:"omitempty,oneof=Pending Triaged Resolved Closed"`
	// 保留字段（向后兼容）
	Title       string `json:"title" binding:"omitempty,max=255"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"omitempty,max=50"`
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
		SelfAssessment:      req.SelfAssessment,
		VulnerabilityURL:    req.VulnerabilityURL,
		VulnerabilityDetail: req.VulnerabilityDetail,
		AttachmentURL:       req.AttachmentURL,
		Severity:            req.Severity,
		AuthorID:            userID.(uint),
		// 保留字段（向后兼容）
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
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
func (h *ReportHandler) ListHandler(c *gin.Context) {
	// 获取 query 参数 ?page=1&page_size=10
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	reports, total, err := h.Service.ListReports(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reports,
		"total": total,
		"page":  page,
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
		SelfAssessment:      req.SelfAssessment,
		VulnerabilityURL:    req.VulnerabilityURL,
		VulnerabilityDetail: req.VulnerabilityDetail,
		AttachmentURL:       req.AttachmentURL,
		Severity:            req.Severity,
		Status:              req.Status,
		// 保留字段（向后兼容）
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}
