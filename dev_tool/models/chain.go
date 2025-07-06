package models

import (
	"gorm.io/gorm"
)

// ChainType 链类型枚举
type ChainType string

const (
	MainNet ChainType = "mainNet"
	TestNet ChainType = "testNet"
	RegTest ChainType = "regTest"
)

// Chain 链信息模型
type Chain struct {
	gorm.Model
	Name         string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	BroadcastURL string    `gorm:"size:255;not null" json:"broadcast_url"`
	UtxoURL      string    `gorm:"size:255" json:"utxo_url"`
	ChainType    ChainType `gorm:"size:20;not null" json:"chain_type"`
}

// TableName 指定表名
func (Chain) TableName() string {
	return "chains"
}
