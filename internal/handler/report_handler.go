package handler

import (
	"bug-bounty-lite/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReportHandler struct {
	Service domain.ReportService
}

func NewReportHandler(s domain.ReportService) *ReportHandler {
	return &ReportHandler{Service: s}
}

// CreateHandler 提交漏洞
func (h *ReportHandler) CreateHandler(c *gin.Context) {
	var report domain.Report
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// [临时逻辑] 假设当前登录用户 ID 是 1
	// 后续我们会从 Context 中获取真实的用户 ID
	report.AuthorID = 1

	if err := h.Service.SubmitReport(&report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": report})
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
	id, _ := strconv.Atoi(idStr)

	report, err := h.Service.GetReport(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}