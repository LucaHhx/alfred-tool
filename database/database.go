package database

import (
	"log"
	"os"
	"path/filepath"

	"sshd/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("无法获取用户主目录:", err)
	}

	dbPath := filepath.Join(homeDir, ".sshd", "connections.db")

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal("无法创建数据库目录:", err)
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("无法连接到数据库:", err)
	}

	// 自动迁移
	if err := DB.AutoMigrate(&models.SSHConnection{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	log.Println("数据库初始化成功")
}

func GetDB() *gorm.DB {
	return DB
}
