package address

type CreateAddressRequest struct {
	Type    string `json:"type" binding:"required,oneof=taproot segwit"`
	ChainID uint   `json:"chain_id" binding:"required"`
}

type AddressResponse struct {
	ID         uint   `json:"id"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key"`
	Type       string `json:"type"`
	ChainID    uint   `json:"chain_id"`
	CreatedAt  string `json:"created_at"`
}
