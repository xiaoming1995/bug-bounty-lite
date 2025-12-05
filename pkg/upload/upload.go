package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// MaxFileSize 最大文件大小 10MB
	MaxFileSize = 10 * 1024 * 1024
	// UploadDir 上传目录
	UploadDir = "uploads/reports"
)

// AllowedMimeTypes 允许的 MIME 类型
var AllowedMimeTypes = map[string]bool{
	"application/pdf":                          true,
	"image/jpeg":                               true,
	"image/jpg":                                true,
	"image/png":                                true,
	"image/gif":                                true,
	"application/msword":                       true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// UploadFile 上传单个文件
func UploadFile(fileHeader *multipart.FileHeader, baseURL string) (*UploadResult, error) {
	// 1. 验证文件大小
	if fileHeader.Size > MaxFileSize {
		return nil, fmt.Errorf("文件大小超过限制（最大10MB）")
	}

	// 2. 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 3. 验证 MIME 类型
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}
	file.Seek(0, 0) // 重置文件指针

	mimeType := http.DetectContentType(buffer)
	if !AllowedMimeTypes[mimeType] {
		return nil, fmt.Errorf("不支持的文件类型: %s", mimeType)
	}

	// 4. 生成文件名和路径
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	
	// 生成唯一文件名
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	
	// 创建目录
	dir := filepath.Join(UploadDir, year, month)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	// 完整路径
	fullPath := filepath.Join(dir, filename)

	// 5. 保存文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 6. 生成访问 URL
	url := fmt.Sprintf("%s/%s/%s/%s", strings.TrimSuffix(baseURL, "/"), UploadDir, year, month, filename)

	return &UploadResult{
		URL:      url,
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: mimeType,
	}, nil
}

