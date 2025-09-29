package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"apps_v1/models"
	"os"
)

var DB *gorm.DB
func ConnectDB(){
	env := os.Getenv("APP_ENV") // misal APP_ENV=local / dev / prod
	host:=os.Getenv("DB_HOST")
	if env == "local" || env == "development" {
		host=os.Getenv("DB_HOST")
	}
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
        host,
        os.Getenv("DB_PORT"),
    )
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
	}

	DB = connection 
	fmt.Println("Connected to PostgreSQL!")
	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Role{})
	connection.AutoMigrate(&models.AccessRight{})
	connection.AutoMigrate(&models.PageApps{})
}
