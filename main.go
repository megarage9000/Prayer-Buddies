package main

import (
	"database/sql"
	"fmt"
	"net/http"

	// Importing the sql drivers
	_ "github.com/lib/pq"
)

// Global config variable
var config Config

func main() {

	var err error

	config, err = LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	// opening the db connection
	_, err = sql.Open("postgres", config.DbURL)

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: serverMux,
	}

	serverMux.HandleFunc("/", helloWorld)

	fmt.Printf("Loading on localhost:%s", config.Port)
	server.ListenAndServe()
}
