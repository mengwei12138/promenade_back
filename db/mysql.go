package db

import (
	"promanage/backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "goMysql:Songge@123@tcp(82.157.57.128:3306)/gomysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	// 自动迁移Requirement表
	db.AutoMigrate(&models.Requirement{})
	return nil
}
