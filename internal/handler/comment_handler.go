package handler

import (
	"bug-bounty-lite/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	service domain.CommentService
}

// NewCommentHandler 创建评论处理器实例
func NewCommentHandler(service domain.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CreateComment 创建评论
// POST /api/reports/:id/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// 获取报告ID
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报告ID"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求体
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	// 创建评论
	comment, err := h.service.CreateComment(uint(reportID), userID.(uint), req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"data":    comment,
	})
}

// ListComments 获取评论列表
// GET /api/reports/:id/comments
func (h *CommentHandler) ListComments(c *gin.Context) {
	// 获取报告ID
	reportIDStr := c.Param("id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的报告ID"})
		return
	}

	// 获取评论列表
	comments, err := h.service.GetReportComments(uint(reportID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": len(comments),
	})
}

// DeleteComment 删除评论
// DELETE /api/reports/:id/comments/:commentId
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// 获取评论ID
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	// 获取当前用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userRole, _ := c.Get("userRole")
	roleStr := ""
	if userRole != nil {
		roleStr = userRole.(string)
	}

	// 删除评论
	if err := h.service.DeleteComment(uint(commentID), userID.(uint), roleStr); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
