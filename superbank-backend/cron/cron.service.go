package cron

import (
	"bank-backend/config"
	"bank-backend/database"
	"bank-backend/module/deposit"
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func InitCron() {
	c := cron.New()
	db := config.DB

	depositService := deposit.NewDepositContainer().Service

	c.AddFunc("0 1 * * *", func() {
		fmt.Println("Running autowithdraw deposit...")

		err := database.WithTransaction(db, func(tx *gorm.DB) error {
			ctxWithTx := database.NewContext(context.Background(), tx)

			err := depositService.ProcessMatureDeposits(ctxWithTx)
			if err != nil {
				fmt.Println("Error in ProcessMatureDeposits:", err.Error())
			}

			fmt.Println("Successfully processed mature deposits")
			return nil
		})

		if err != nil {
			fmt.Println("Transaction failed:", err.Error())
		}
	})

	c.Start()
}
