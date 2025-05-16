package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (config *Config) CreateUser(resp http.ResponseWriter, req *http.Request) {

	// Decoding the response
	decoder := json.NewDecoder(req.Body)
	result := User{}

	err := decoder.Decode(&result)
	if err != nil {
		message := "ERROR: unable to decode request"
		fmt.Printf(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
	}

}
