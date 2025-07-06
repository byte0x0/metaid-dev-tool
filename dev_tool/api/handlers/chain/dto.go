package chain

// CreateChainRequest 创建链的请求结构
type CreateChainRequest struct {
	Name         string `json:"name" binding:"required"`
	BroadcastURL string `json:"broadcastUrl" binding:"required,url"`
	UtxoURL      string `json:"utxoUrl" binding:"omitempty,url"`
	ChainType    string `json:"chainType" binding:"required,oneof=mainNet testNet regTest"`
}

// UpdateChainRequest 更新链的请求结构
type UpdateChainRequest struct {
	Name         string `json:"name,omitempty"`
	BroadcastURL string `json:"broadcastUrl,omitempty" binding:"omitempty,url"`
	UtxoURL      string `json:"utxoUrl" binding:"omitempty,url"`
	ChainType    string `json:"chainType,omitempty" binding:"omitempty,oneof=mainNet testNet regTest"`
}
