package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	app, err := CreateApp()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
