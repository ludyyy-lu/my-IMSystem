package model

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "imuser:im123456@tcp(localhost:3307)/im?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移建表
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	fmt.Println("✅ MySQL connected and User table migrated.")
}
