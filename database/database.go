package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)


func ConnectDB() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	dsn := os.Getenv("DATABASE_URL")
	Connect, err := gorm.Open(postgres.Open(dsn)) 

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db = Connect
}


func GetDB() *gorm.DB{
	return db
}

