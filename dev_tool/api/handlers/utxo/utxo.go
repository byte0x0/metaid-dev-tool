package utxo

import (
	"dev_tool/api/response"

	"github.com/gin-gonic/gin"
)

// GetUTXO 获取UTXO信息
func GetUTXO(c *gin.Context) {
	id := c.Param("id")

	data := gin.H{
		"id":      id,
		"amount":  "1.23",
		"address": "addr_xxx",
		"spent":   false,
	}

	response.Success(c, data)
}

// GetUTXOsByAddress 获取地址的UTXO列表
func GetUTXOsByAddress(c *gin.Context) {
	address := c.Param("address")

	data := gin.H{
		"address": address,
		"utxos": []gin.H{
			{
				"id":     "utxo1",
				"amount": "1.23",
				"spent":  false,
			},
			{
				"id":     "utxo2",
				"amount": "2.34",
				"spent":  false,
			},
		},
	}

	response.Success(c, data)
}
