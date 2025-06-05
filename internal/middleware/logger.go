package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求参数日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后需要重置
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 获取请求参数
		queryParams := c.Request.URL.Query()
		pathParams := c.Params

		// 记录请求信息
		requestInfo := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"query":       queryParams,
			"path_params": pathParams,
			"headers":     c.Request.Header,
		}

		// 如果有请求体，尝试解析为 JSON
		if len(requestBody) > 0 {
			var jsonBody interface{}
			if err := json.Unmarshal(requestBody, &jsonBody); err == nil {
				requestInfo["body"] = jsonBody
			} else {
				requestInfo["body"] = string(requestBody)
			}
		}

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 获取响应状态
		statusCode := c.Writer.Status()

		// 记录响应信息
		responseInfo := map[string]interface{}{
			"status_code": statusCode,
			"latency":     latency.String(),
		}

		// 获取响应体
		if c.Writer.Written() {
			responseInfo["response"] = c.Writer.Header().Get("Content-Type")
		}

		// 打印日志
		if statusCode >= 400 {
			// 错误请求使用红色标记
			c.Error(fmt.Errorf("请求失败: %v", responseInfo))
		}

		// 打印请求和响应信息
		requestJSON, _ := json.MarshalIndent(requestInfo, "", "  ")
		responseJSON, _ := json.MarshalIndent(responseInfo, "", "  ")

		gin.DefaultWriter.Write([]byte("\n=== 请求信息 ===\n"))
		gin.DefaultWriter.Write(requestJSON)
		gin.DefaultWriter.Write([]byte("\n=== 响应信息 ===\n"))
		gin.DefaultWriter.Write(responseJSON)
		gin.DefaultWriter.Write([]byte("\n================\n"))
	}
}
