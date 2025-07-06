package routes

import (
	"dev_tool/api/handlers/inscribe"

	"github.com/gin-gonic/gin"
)

func RegisterInscribeRoutes(r *gin.RouterGroup) {
	inscribeGroup := r.Group("/inscribes")
	{
		inscribeGroup.POST("", inscribe.CreateInscribe)
		inscribeGroup.GET("", inscribe.ListInscribes)
		inscribeGroup.GET("/:id", inscribe.GetInscribe)
	}
}
