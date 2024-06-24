package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)


func ConnectDB() {
	dsn := "postgresql://shopdb_owner:TZOEtPp4N1uj@ep-steep-moon-a5sk29u2.us-east-2.aws.neon.tech/shopdb?sslmode=require"
	Connect, err := gorm.Open(postgres.Open(dsn)) 

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db = Connect
}


func GetDB() *gorm.DB{
	return db
}