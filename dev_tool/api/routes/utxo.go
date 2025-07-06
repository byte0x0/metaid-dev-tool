package routes

import (
	"dev_tool/api/handlers/proxy"
	"dev_tool/api/handlers/utxo"

	"github.com/gin-gonic/gin"
)

func RegisterUTXORoutes(r *gin.RouterGroup) {
	utxoGroup := r.Group("/utxo")
	{
		utxoGroup.GET("/:id", utxo.GetUTXO)
		utxoGroup.GET("/address/:address", utxo.GetUTXOsByAddress)
	}

	// 添加代理路由
	proxyGroup := r.Group("/proxy")
	{
		proxyGroup.POST("/utxo", proxy.ProxyUTXO)
	}
}
