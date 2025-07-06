package middleware

import (
	"dev_tool/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 处理跨域请求中间件
func Cors() gin.HandlerFunc {
	corsConfig := config.GlobalConfig.Cors

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查来源是否在允许列表中
		allowOrigin := "*"
		if len(corsConfig.AllowOrigins) > 0 && corsConfig.AllowOrigins[0] != "*" {
			for _, allowed := range corsConfig.AllowOrigins {
				if allowed == origin {
					allowOrigin = origin
					break
				}
			}
		}

		// 设置 CORS 头
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowMethods, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowHeaders, ", "))
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", string(corsConfig.MaxAge))

		if corsConfig.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
