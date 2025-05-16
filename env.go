package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DbURL string
}

func LoadConfig() (Config, error) {

	err := godotenv.Load()

	if err != nil {
		return Config{}, fmt.Errorf("Unable to load environment variables")
	}

	portNumber := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	return Config{
			Port:  portNumber,
			DbURL: dbURL,
		},
		nil
}
