package routes

import (
	"dev_tool/api/handlers/broadcast"

	"github.com/gin-gonic/gin"
)

func RegisterBroadcastRoutes(r *gin.RouterGroup) {
	broadcastGroup := r.Group("/broadcast")
	{
		broadcastGroup.POST("", broadcast.ProxyBroadcast)
	}

}
