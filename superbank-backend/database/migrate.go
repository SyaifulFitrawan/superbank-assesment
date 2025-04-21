package database

import (
	"bank-backend/config"
	"bank-backend/model"
	"bank-backend/utils"
	"fmt"
)

func MigrateDB() {
	logger := utils.NewLogger()
	db := config.DB

	if err := db.AutoMigrate(
		&model.User{},
		&model.Customer{},
		&model.Deposit{},
		&model.Pocket{},
	); err != nil {
		logger.Error("Failed to migrate database;", err.Error())
		return
	}

	fmt.Println("Database migrated successfully!")
}
