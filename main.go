package main

import (
	"fmt"
	"net/http"
)

func main() {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: serverMux,
	}

	serverMux.HandleFunc("/", helloWorld)

	fmt.Printf("Loading on localhost:%s", config.Port)
	server.ListenAndServe()
}
