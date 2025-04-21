package seeders

import (
	"bank-backend/config"
	seeder_data "bank-backend/seeders/data"
	"bank-backend/utils"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

func UserSeeds() error {
	var users = seeder_data.DummyUser
	logger := utils.NewLogger()

	for i := range users {
		if users[i].ID == uuid.Nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(users[i].Password), bcrypt.DefaultCost)
			if err != nil {
				continue
			}

			users[i].ID = uuid.NewV4()
			users[i].Password = string(hashedPassword)
		}
	}

	if err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "password", "updated_at"}),
	}).Create(&users).Error; err != nil {
		logger.Error("Failed to upsert users", err.Error())
		return err
	}

	logger.Log("✅ User seeding success")
	return nil
}
