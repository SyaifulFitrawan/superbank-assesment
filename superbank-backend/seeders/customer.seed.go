package seeders

import (
	"bank-backend/config"
	"bank-backend/model"
	seeder_data "bank-backend/seeders/data"
	"bank-backend/utils"
	"time"

	"github.com/bxcodec/faker/v4"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm/clause"
)

func CustomerSeeds() error {
	var customers = seeder_data.DummyCustomer
	logger := utils.NewLogger()

	for i := range customers {
		if customers[i].ID == uuid.Nil {
			customers[i].ID = uuid.NewV4()
			customers[i].Balance = utils.RandomRoundedAmount()
		}
	}

	if err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_number"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
	}).Create(&customers).Error; err != nil {
		logger.Error("Failed to upsert customers", err.Error())
		return err
	}

	for _, customer := range customers {
		numDeposits := utils.RandomInt(1, 3)
		for i := 0; i < numDeposits; i++ {
			startDate := utils.RandomDate("2006-01-02")
			term := utils.RandomInt(6, 36)
			maturityDate := startDate.AddDate(0, term, 0)

			deposit := model.Deposit{
				ID:           uuid.NewV4(),
				CustomerID:   customer.ID,
				Amount:       utils.RandomRoundedAmount(),
				InterestRate: utils.RandomFloat(2.5, 8.0),
				TermMonths:   term,
				StartDate:    startDate,
				MaturityDate: maturityDate,
				IsWithdrawn:  false,
				Note:         faker.Sentence(),
			}

			if err := config.DB.Create(&deposit).Error; err != nil {
				logger.Error("Failed to create deposit", err.Error())
				continue
			}
		}

		numPockets := utils.RandomInt(0, 4)
		for i := 0; i < numPockets; i++ {
			var targetAmount *float64
			var targetDate *time.Time
			if utils.RandomBool() {
				amount := utils.RandomRoundedAmount()
				targetAmount = &amount
			}
			if utils.RandomBool() {
				tDate := time.Now().AddDate(0, utils.RandomInt(3, 12), 0)
				targetDate = &tDate
			}

			pocket := model.Pocket{
				ID:           uuid.NewV4(),
				CustomerID:   customer.ID,
				Name:         faker.Word() + " Goal",
				Balance:      utils.RandomRoundedAmount(),
				TargetAmount: targetAmount,
				TargetDate:   targetDate,
				IsActive:     true,
			}

			if err := config.DB.Create(&pocket).Error; err != nil {
				logger.Error("Failed to create pocket", err.Error())
				continue
			}
		}
	}

	logger.Log("✅ Customer, Deposit, and Pocket seeding success")
	return nil
}
