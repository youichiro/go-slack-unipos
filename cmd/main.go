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

	db, err := repository.InitDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	r := router.SetupRouter(db)
	err = r.Run(":8080")
	if err != nil {
		panic(err.Error())
	}
}
