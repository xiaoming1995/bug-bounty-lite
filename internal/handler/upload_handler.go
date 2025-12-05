package handler

import (
	"fmt"
	"bug-bounty-lite/pkg/response"
	"bug-bounty-lite/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadHandler 文件上传处理器
type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadFileHandler 上传单个文件
// POST /api/v1/upload
func (h *UploadHandler) UploadFileHandler(c *gin.Context) {
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择要上传的文件")
		return
	}

	// 获取 base URL（用于生成文件访问 URL）
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	// 上传文件
	result, err := upload.UploadFile(file, baseURL)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, result)
}

