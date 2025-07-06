package routes

import (
	"dev_tool/api/handlers/chain"

	"github.com/gin-gonic/gin"
)

func RegisterChainRoutes(r *gin.RouterGroup) {
	chainGroup := r.Group("/chains")
	{
		chainGroup.POST("", chain.CreateChain)       // 创建链
		chainGroup.GET("", chain.ListChains)         // 获取链列表
		chainGroup.GET("/:id", chain.GetChain)       // 获取单个链
		chainGroup.PUT("/:id", chain.UpdateChain)    // 更新链
		chainGroup.DELETE("/:id", chain.DeleteChain) // 删除链
	}
}
