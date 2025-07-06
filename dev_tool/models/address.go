package models

import (
	"gorm.io/gorm"
)

type AddressType string

const (
	Taproot AddressType = "taproot"
	Segwit  AddressType = "segwit"
)

type Address struct {
	gorm.Model
	Address    string      `gorm:"size:100;not null;uniqueIndex" json:"address"`
	PrivateKey string      `gorm:"size:100;not null" json:"private_key"`
	Type       AddressType `gorm:"size:20;not null" json:"type"`
	ChainID    uint        `gorm:"not null" json:"chain_id"` // 关联的链ID
	Chain      Chain       `gorm:"foreignKey:ChainID" json:"chain"`
}

func (Address) TableName() string {
	return "addresses"
}
