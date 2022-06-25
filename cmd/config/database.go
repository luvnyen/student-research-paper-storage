package config

import (
	"fmt"

	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", "root", "", "student-research-paper")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Student{}, &models.Paper{})
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic(err)
	}
	dbSQL.Close()
}
