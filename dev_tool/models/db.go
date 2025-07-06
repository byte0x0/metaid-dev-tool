package models

import (
	"log"

	"dev_tool/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "modernc.org/sqlite"
)

var DB *gorm.DB

func InitDB() {
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(sqlite.Open(config.GlobalConfig.Database.SQLite.Path), gormConfig)
	if err != nil {
		log.Fatalf("连接数据库失败 (%s): %v", config.GlobalConfig.Database.SQLite.Path, err)
	}

	log.Println("数据库连接成功")

	// 自动迁移数据库结构
	err = autoMigrate()
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}

func autoMigrate() error {
	return DB.AutoMigrate(&Chain{}, &Address{})
}
