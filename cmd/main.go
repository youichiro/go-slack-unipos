package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/youichiro/go-slack-my-unipos/internal/repository"
	"github.com/youichiro/go-slack-my-unipos/internal/router"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbRepo := &repository.PostgresRepository{}
	err = dbRepo.Connect()
	if err != nil {
		panic(err.Error())
	}
	defer dbRepo.Close()

	r := router.SetupRouter()
	err = r.Run(":8080")
	if err != nil {
		panic(err.Error())
	}
}
