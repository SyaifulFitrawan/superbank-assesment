package database

import (
	"bank-backend/seeders"
	"fmt"
)

func SeedAllData() {
	fmt.Println("Seeding process started...")

	seeders.UserSeeds()
	seeders.CustomerSeeds()

	fmt.Println("All seeding processes completed!")
}
