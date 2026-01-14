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
	Category    string `json:"category"`
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Category    string `json:"category"`
}

// ReviewRequest 审核请求
type ReviewRequest struct {
	Approved     bool   `json:"approved"`
	RejectReason string `json:"reject_reason"`
}

// SetFeaturedRequest 设置精选请求
type SetFeaturedRequest struct {
	Featured bool `json:"featured"`
}

// CreateArticle 创建文章
// POST /api/v1/articles
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
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

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	article, err := h.service.CreateArticle(userID.(uint), roleStr, req.Title, req.Description, req.Content, req.Category)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := "文章创建成功，等待审核"
	if roleStr == "admin" {
		message = "文章发布成功"
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": message,
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

	article, err := h.service.UpdateArticle(uint(articleID), userID.(uint), req.Title, req.Description, req.Content, req.Category)
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

	// 获取客户端 IP
	clientIP := c.ClientIP()

	article, err := h.service.GetArticle(uint(articleID), incrementView, clientIP)
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

// GetFeaturedArticles 获取精选文章
// GET /api/v1/articles/public/featured
func (h *ArticleHandler) GetFeaturedArticles(c *gin.Context) {
	limit := 3 // 默认3篇
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	articles, err := h.service.GetFeaturedArticles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取精选文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
	})
}

// GetHotArticles 获取热门文章
// GET /api/v1/articles/public/hot
func (h *ArticleHandler) GetHotArticles(c *gin.Context) {
	limit := 3 // 默认3篇
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	articles, err := h.service.GetHotArticles(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": len(articles),
	})
}

// SetFeatured 设置精选状态（管理员）
// PUT /api/v1/admin/articles/:id/featured
func (h *ArticleHandler) SetFeatured(c *gin.Context) {
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

	var req SetFeaturedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效"})
		return
	}

	if err := h.service.SetFeatured(uint(articleID), req.Featured); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := "已取消精选"
	if req.Featured {
		message = "已设为精选"
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
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
