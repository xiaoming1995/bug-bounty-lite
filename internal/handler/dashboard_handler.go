package handler

import (
	"bug-bounty-lite/internal/domain"
	"bug-bounty-lite/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	Service domain.DashboardService
}

func NewDashboardHandler(s domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{Service: s}
}

// getUserInfo 从 Context 获取用户 ID 和角色
func getUserInfo(c *gin.Context) (uint, string) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, ""
	}
	userID := userIDVal.(uint)

	roleVal, exists := c.Get("role")
	userRole := "whitehat" // 默认角色
	if exists && roleVal != nil {
		userRole = roleVal.(string)
	}

	return userID, userRole
}

// GetStatistics 获取漏洞统计数据
// GET /api/v1/dashboard/statistics
func (h *DashboardHandler) GetStatistics(c *gin.Context) {
	userID, userRole := getUserInfo(c)

	stats, err := h.Service.GetStatistics(userID, userRole)
	if err != nil {
		response.Error(c, 500, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTrend 获取漏洞趋势数据
// GET /api/v1/dashboard/trend?period=month
func (h *DashboardHandler) GetTrend(c *gin.Context) {
	period := c.DefaultQuery("period", "month")
	userID, userRole := getUserInfo(c)

	trend, err := h.Service.GetTrend(period, userID, userRole)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, trend)
}

// GetReports 获取漏洞列表
// GET /api/v1/dashboard/reports?type=pending&limit=6
func (h *DashboardHandler) GetReports(c *gin.Context) {
	reportType := c.DefaultQuery("type", "pending")
	limitStr := c.DefaultQuery("limit", "6")
	userID, userRole := getUserInfo(c)

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 6
	}

	reports, total, err := h.Service.GetReportsByType(reportType, limit, userID, userRole)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  reports,
		"total": total,
	})
}
