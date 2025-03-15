package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jasonzheng/carrag/utils"
)

// Logger 创建一个Gin中间件，用于记录HTTP请求日志
func Logger(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方法
		method := c.Request.Method
		// 请求路由
		path := c.Request.URL.Path
		// 状态码
		statusCode := c.Writer.Status()
		// 客户端IP
		clientIP := c.ClientIP()

		// 记录日志
		logger.Info("%s | %d | %s | %s | %v", method, statusCode, clientIP, path, latency)
	}
}
