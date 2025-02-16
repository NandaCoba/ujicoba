package database

import (
	"fmt"
	"worker/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed connect DB")
	}

	DB = db

	db.AutoMigrate(models.Karyawan{})
	fmt.Println("Success connect DB")
}
