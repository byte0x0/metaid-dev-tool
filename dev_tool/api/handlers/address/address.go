package address

import (
	"dev_tool/api/response"
	"dev_tool/chain/tool"
	"dev_tool/models"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getNetParams 根据链类型获取对应的网络参数
func getNetParams(chainType models.ChainType) *chaincfg.Params {
	switch chainType {
	case models.MainNet:
		return &chaincfg.MainNetParams
	case models.TestNet:
		return &chaincfg.TestNet3Params
	case models.RegTest:
		return &chaincfg.RegressionNetParams
	default:
		return &chaincfg.TestNet3Params
	}
}

// CreateAddress 创建新地址
func CreateAddress(c *gin.Context) {
	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "无效的请求参数")
		return
	}

	// 查找关联的链信息
	var chain models.Chain
	if err := models.DB.First(&chain, req.ChainID).Error; err != nil {
		response.Error(c, "指定的链不存在")
		return
	}

	// 获取网络参数
	netParams := getNetParams(chain.ChainType)

	// 根据类型创建地址
	var privateKey, addr string
	var err error

	switch models.AddressType(req.Type) {
	case models.Taproot:
		privateKey, addr, err = tool.CreateTaprootKey(netParams)
	case models.Segwit:
		privateKey, addr, err = tool.CreateSegwitKey(netParams)
	default:
		response.Error(c, "不支持的地址类型")
		return
	}

	if err != nil {
		response.Error(c, "创建地址失败: "+err.Error())
		return
	}

	// 保存到数据库
	address := models.Address{
		Address:    addr,
		PrivateKey: privateKey,
		Type:       models.AddressType(req.Type),
		ChainID:    req.ChainID,
	}

	if err := models.DB.Create(&address).Error; err != nil {
		response.Error(c, "保存地址失败")
		return
	}

	response.Success(c, address)
}

// GetAddress 获取单个地址信息
func GetAddress(c *gin.Context) {
	id := c.Param("id")

	var address models.Address
	if err := models.DB.Preload("Chain").First(&address, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Error(c, "地址不存在")
			return
		}
		response.Error(c, "获取地址信息失败")
		return
	}

	response.Success(c, address)
}

// ListAddresses 获取地址列表
func ListAddresses(c *gin.Context) {
	var addresses []models.Address

	query := models.DB.Preload("Chain")

	// 支持按类型和链ID筛选
	if addressType := c.Query("type"); addressType != "" {
		query = query.Where("type = ?", addressType)
	}
	if chainID := c.Query("chain_id"); chainID != "" {
		query = query.Where("chain_id = ?", chainID)
	}

	if err := query.Find(&addresses).Error; err != nil {
		response.Error(c, "获取地址列表失败")
		return
	}

	response.Success(c, addresses)
}

// DeleteAddress 删除地址
func DeleteAddress(c *gin.Context) {
	id := c.Param("id")

	if err := models.DB.Delete(&models.Address{}, id).Error; err != nil {
		response.Error(c, "删除地址失败")
		return
	}

	response.Success(c, gin.H{"id": id})
}
