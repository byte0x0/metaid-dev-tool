package middleware

import (
	"dev_tool/api/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 处理panic错误
				response.Error(c, "Internal Server Error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
