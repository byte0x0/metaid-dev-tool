package chain

import (
	"dev_tool/api/response"
	"dev_tool/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateChain 创建新的链
func CreateChain(c *gin.Context) {
	var req CreateChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		response.Error(c, "无效的请求参数")
		return
	}

	chain := models.Chain{
		Name:         req.Name,
		BroadcastURL: req.BroadcastURL,
		UtxoURL:      req.UtxoURL,
		ChainType:    models.ChainType(req.ChainType),
	}

	if err := models.DB.Create(&chain).Error; err != nil {
		response.Error(c, "创建链失败")
		return
	}

	response.Success(c, chain)
}

// GetChain 获取单个链信息
func GetChain(c *gin.Context) {
	id := c.Param("id")

	var chain models.Chain
	if err := models.DB.First(&chain, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "链不存在")
			return
		}
		response.Error(c, "获取链信息失败")
		return
	}

	response.Success(c, chain)
}

// ListChains 获取所有链列表
func ListChains(c *gin.Context) {
	var chains []models.Chain

	// 支持按类型筛选
	chainType := c.Query("chain_type")
	db := models.DB
	if chainType != "" {
		db = db.Where("chain_type = ?", chainType)
	}

	if err := db.Find(&chains).Error; err != nil {
		response.Error(c, "获取链列表失败")
		return
	}

	response.Success(c, chains)
}

// UpdateChain 更新链信息
func UpdateChain(c *gin.Context) {
	id := c.Param("id")

	var req UpdateChainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "无效的请求参数")
		return
	}

	var chain models.Chain
	if err := models.DB.First(&chain, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "链不存在")
			return
		}
		response.Error(c, "获取链信息失败")
		return
	}

	// 只更新提供的字段
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.BroadcastURL != "" {
		updates["broadcast_url"] = req.BroadcastURL
	}
	// 直接更新 UtxoURL，允许空字符串
	updates["utxo_url"] = req.UtxoURL
	if req.ChainType != "" {
		updates["chain_type"] = req.ChainType
	}

	if err := models.DB.Model(&chain).Updates(updates).Error; err != nil {
		response.Error(c, "更新链信息失败")
		return
	}

	// 重新获取更新后的完整数据
	if err := models.DB.First(&chain, id).Error; err != nil {
		response.Error(c, "获取更新后的链信息失败")
		return
	}

	response.Success(c, chain)
}

// DeleteChain 删除链
func DeleteChain(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Chain{}, id).Error; err != nil {
		response.Error(c, "删除链失败")
		return
	}

	response.Success(c, gin.H{"id": id})
}
