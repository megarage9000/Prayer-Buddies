package main

import (
	"fmt"
	"net/http"
	"time"

	// Importing the sql drivers
	_ "github.com/lib/pq"
)

func main() {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.Port),
		Handler:           serverMux,
		ReadHeaderTimeout: time.Minute,
	}

	serverMux.HandleFunc("POST /api/users", config.CreateUser)
	serverMux.HandleFunc("POST /api/login", config.LoginUser)

	fmt.Printf("Loading on localhost:%s", config.Port)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}
}
