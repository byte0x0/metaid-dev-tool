package api

import (
	"dev_tool/api/middleware"
	"dev_tool/api/response"
	"dev_tool/api/routes"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 注册自定义验证器
func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("url_if_exists", func(fl validator.FieldLevel) bool {
		// 如果字段是指针类型且为 nil，返回 true（验证通过）
		if ptr, ok := fl.Field().Interface().(*string); ok {
			if ptr == nil {
				return true
			}
			// 如果不为 nil，则验证 URL 格式
			if v.Var(*ptr, "url") == nil {
				return true
			}
		}
		return false
	})
}

// SetupRouter 初始化所有路由
func SetupRouter(r *gin.Engine) {
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		registerCustomValidators(v)
	}

	// 添加 CORS 中间件
	r.Use(middleware.Cors())

	// 添加错误处理中间件
	r.Use(middleware.ErrorHandler())

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 健康检查路由
		v1.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{"message": "pong"})
		})

		// 注册各模块路由
		routes.RegisterAddressRoutes(v1)
		routes.RegisterChainRoutes(v1)
		routes.RegisterUTXORoutes(v1)
		routes.RegisterInscribeRoutes(v1)
		routes.RegisterBroadcastRoutes(v1)
	}
}
