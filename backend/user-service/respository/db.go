package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"user-services/models"
)

var DB *gorm.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Check valid port and convert to number
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return fmt.Errorf("invalid DB_PORT: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Set DB connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("failed to connect to DB: %v", err)
		return err
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Error during AutoMigrate: %v", err)
	}

	log.Println("Successfully connected to DB and set up")
	return nil
}
