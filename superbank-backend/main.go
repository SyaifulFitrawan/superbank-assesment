package main

import (
	"bank-backend/config"
	"bank-backend/cron"
	"bank-backend/database"
	"bank-backend/router"
	"bank-backend/utils"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	app := fiber.New()
	app.Use(cors.New())
	config.ConnectDatabase()
	utils.InitValidator()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			database.MigrateDB()
			return
		case "seed":
			database.SeedAllData()
			return
		}
	}

	cron.InitCron()
	router.SetupRoutes(app)

	port := fmt.Sprintf(":%s", config.Env("PORT"))
	err := app.Listen(port)
	if err != nil {
		log.Fatalf("fiber.Listen failed %s", err)
	}
}
