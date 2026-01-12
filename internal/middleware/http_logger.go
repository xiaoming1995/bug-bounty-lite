package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter 包装 gin.ResponseWriter 以捕获响应体内容
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// HttpLogger 返回一个高级日志中间件，记录完整的 HTTP 交互过程
func HttpLogger() gin.HandlerFunc {
	// 确保日志目录存在
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.Mkdir(logDir, 0755)
	}

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// --- 捕获请求体 (Request Body) ---
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 把原 Body 写回去，方便后续读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// --- 捕获响应体 (Response Body) ---
		writer := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer

		// 处理请求
		c.Next()

		// --- 记录日志 ---
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 格式化日志
		logEntry := fmt.Sprintf(
			"[%s] %s | %d | %v | %s | %s\n"+
				"  [Request Headers]: %v\n"+
				"  [Request Body]: %s\n"+
				"  [Response Body]: %s\n"+
				"----------------------------------------------------------------------\n",
			endTime.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.URL.RequestURI(),
			c.Request.Header,
			string(requestBody),
			writer.body.String(),
		)

		// 异步写入文件，按日期切分
		go func(entry string) {
			fileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
			f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("[Error] Failed to open log file: %v", err)
				return
			}
			defer f.Close()
			_, _ = f.WriteString(entry)
		}(logEntry)
	}
}
