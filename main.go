package main

import (
	"net/http"
)

func main() {

	serverMux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	serverMux.HandleFunc("/", helloWorld)
	server.ListenAndServe()
}

