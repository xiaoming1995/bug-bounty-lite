package handler

import (
	"bug-bounty-lite/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleLikeCommentHandler struct {
	service service.ArticleLikeCommentService
}

// NewArticleLikeCommentHandler 创建点赞评论处理器
func NewArticleLikeCommentHandler(service service.ArticleLikeCommentService) *ArticleLikeCommentHandler {
	return &ArticleLikeCommentHandler{service: service}
}

// ToggleLike 切换点赞状态
// POST /api/v1/articles/:id/like
func (h *ArticleLikeCommentHandler) ToggleLike(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	liked, likeCount, err := h.service.ToggleLike(uint(articleID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "操作成功",
		"liked":      liked,
		"like_count": likeCount,
	})
}

// GetLikeStatus 获取点赞状态
// GET /api/v1/articles/:id/like
func (h *ArticleLikeCommentHandler) GetLikeStatus(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 获取用户ID（可选，未登录用户返回 liked: false）
	var userID uint = 0
	if id, exists := c.Get("userID"); exists {
		userID = id.(uint)
	}

	liked, likeCount, err := h.service.GetLikeStatus(uint(articleID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked":      liked,
		"like_count": likeCount,
	})
}

// AddComment 发表评论
// POST /api/v1/articles/:id/comments
func (h *ArticleLikeCommentHandler) AddComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "评论内容不能为空"})
		return
	}

	comment, err := h.service.AddComment(uint(articleID), userID.(uint), req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论成功",
		"data":    comment,
	})
}

// GetComments 获取评论列表
// GET /api/v1/articles/:id/comments
func (h *ArticleLikeCommentHandler) GetComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	comments, err := h.service.GetComments(uint(articleID))
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
// DELETE /api/v1/articles/:id/comments/:commentId
func (h *ArticleLikeCommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的评论ID"})
		return
	}

	if err := h.service.DeleteComment(uint(commentID), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
