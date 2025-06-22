package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtKey = []byte(os.Getenv("JWT_SECRET"))

func ConnectPostgre() {
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN environment variable is not set")
	}
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %v", err)
	}
}