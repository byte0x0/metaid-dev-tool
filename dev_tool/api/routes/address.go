package routes

import (
	"dev_tool/api/handlers/address"

	"github.com/gin-gonic/gin"
)

func RegisterAddressRoutes(r *gin.RouterGroup) {
	addressGroup := r.Group("/addresses")
	{
		addressGroup.POST("", address.CreateAddress)
		addressGroup.GET("", address.ListAddresses)
		addressGroup.GET("/:id", address.GetAddress)
		addressGroup.DELETE("/:id", address.DeleteAddress)
	}
}
