package db

import (
	"fmt"
	"gadfix/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// connection
func Connection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}
	root := os.Getenv("Root")

	database, err := gorm.Open(mysql.Open(root), &gorm.Config{})
	if err != nil {

		log.Fatal("failed to connect to database", err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.UserAddress{},
		&models.Service{},
		&models.Staff{},
		&models.Booking{},
	)

	if err != nil {
		log.Fatal("failed to auto migrate")
	}

	DB = database
	fmt.Println("database connected successfuly")
}
