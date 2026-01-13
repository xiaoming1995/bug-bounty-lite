package handler

import (
	"bug-bounty-lite/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	service domain.ArticleService
}

// NewArticleHandler 创建文章处理器实例
func NewArticleHandler(service domain.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
}

// ReviewRequest 审核请求
type ReviewRequest struct {
	Approved     bool   `json:"approved"`
	RejectReason string `json:"reject_reason"`
}

// CreateArticle 创建文章
// POST /api/v1/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	article, err := h.service.CreateArticle(userID.(uint), req.Title, req.Description, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功，等待审核",
		"data":    article,
	})
}

// UpdateArticle 更新文章
// PUT /api/v1/articles/:id
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	article, err := h.service.UpdateArticle(uint(articleID), userID.(uint), req.Title, req.Description, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"data":    article,
	})
}

// DeleteArticle 删除文章
// DELETE /api/v1/articles/:id
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
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

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	if err := h.service.DeleteArticle(uint(articleID), userID.(uint), roleStr); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}

// GetArticle 获取文章详情
// GET /api/v1/articles/:id
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 判断是否需要增加浏览量（公开访问时增加）
	incrementView := c.Query("view") == "true"

	article, err := h.service.GetArticle(uint(articleID), incrementView)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": article})
}

// GetMyArticles 获取我的文章列表
// GET /api/v1/articles
func (h *ArticleHandler) GetMyArticles(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	articles, err := h.service.GetMyArticles(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
	})
}

// GetPublishedArticles 获取已发布的文章列表（学习中心）
// GET /api/v1/articles/public
func (h *ArticleHandler) GetPublishedArticles(c *gin.Context) {
	articles, err := h.service.GetPublishedArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
	})
}

// ReviewArticle 审核文章（管理员）
// PUT /api/v1/admin/articles/:id/review
func (h *ArticleHandler) ReviewArticle(c *gin.Context) {
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权执行此操作"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	article, err := h.service.ReviewArticle(uint(articleID), req.Approved, req.RejectReason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "审核完成",
		"data":    article,
	})
}
