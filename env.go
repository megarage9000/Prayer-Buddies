package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

type Config struct {
	Port        string
	Secret      string
	FrontendURL string
	Database    *database.Queries
}

func LoadConfig() (Config, error) {

	err := godotenv.Load()

	if err != nil {
		return Config{}, fmt.Errorf("unable to load environment variables")
	}

	portNumber := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("SECRET")
	frontendURL := os.Getenv("FRONTEND_URL")

	// opening the db connection
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		return Config{}, fmt.Errorf("unable to connect to database")
	}

	return Config{
			Port:        portNumber,
			Secret:      secret,
			Database:    database.New(db),
			FrontendURL: frontendURL,
		},
		nil
}
