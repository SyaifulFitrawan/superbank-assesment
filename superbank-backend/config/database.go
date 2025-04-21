package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := Env("DB_HOST")
	dbPort := Env("DB_PORT")
	dbUser := Env("DB_USER")
	dbPass := Env("DB_PASS")
	dbName := Env("DB_NAME")

	if dbPort == "" || dbPort == "0" {
		dbPort = "5432"
	}

	if dbHost == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("Database environment variables are missing! Check your .env file.")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	var errDB error
	DB, errDB = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		log.Fatal("Database connection error:", errDB)
	}

	fmt.Println("Database connected successfully!")
}
