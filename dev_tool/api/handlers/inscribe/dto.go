package inscribe

// CreateInscribeRequest 创建铭文请求
type CreateInscribeRequest struct {
	ChainID     uint    `json:"chainId" binding:"required"`
	MetaIDFlag  string  `json:"metaIdFlag" binding:"required"`
	Operation   string  `json:"operation" binding:"required"`
	Path        string  `json:"path" binding:"required"`
	Payload     string  `json:"payload" binding:"required"` // 改为 string 类型
	PinOutValue int64   `json:"pinOutValue" binding:"required,min=546"`
	Address     string  `json:"address" binding:"required"`
	FeeRate     int64   `json:"feeRate" binding:"required,min=1"`
	UtxoList    []*Utxo `json:"utxoList" binding:"required,min=1"`
	CommitTx    string  `json:"commitTx,omitempty"`
	RevealTx    string  `json:"revealTx,omitempty"`
}

type InscribeResponse struct {
	CommitTxRaw string `json:"commit_tx_raw"`
	RevealTxRaw string `json:"reveal_tx_raw"`
	CommitTxID  string `json:"commit_tx_id"`
	RevealTxID  string `json:"reveal_tx_id"`
	MinerFee    int64  `json:"miner_fee"`
}
