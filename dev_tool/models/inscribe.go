package models

import (
	"gorm.io/gorm"
)

type InscribeStatus string

const (
	Pending   InscribeStatus = "pending"
	Committed InscribeStatus = "committed"
	Revealed  InscribeStatus = "revealed"
	Failed    InscribeStatus = "failed"
)

type Inscribe struct {
	gorm.Model
	ChainID     uint           `gorm:"not null" json:"chain_id"`
	Chain       Chain          `gorm:"foreignKey:ChainID" json:"chain"`
	MetaIDFlag  string         `gorm:"size:100;not null" json:"meta_id_flag"`
	Path        string         `gorm:"size:255;not null" json:"path"`
	Payload     string         `gorm:"type:text;not null" json:"payload"`
	Address     string         `gorm:"size:100;not null" json:"address"`
	CommitTxID  string         `gorm:"size:64" json:"commit_tx_id"`
	RevealTxID  string         `gorm:"size:64" json:"reveal_tx_id"`
	Status      InscribeStatus `gorm:"size:20;not null" json:"status"`
	PinOutValue int64          `gorm:"not null" json:"pin_out_value"`
	FeeRate     int64          `gorm:"not null" json:"fee_rate"`
	MinerFee    int64          `json:"miner_fee"`
}
