package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Function to return JSON
func RespondJSON(resp http.ResponseWriter, req *http.Request, payload interface{}, statusCode int) {
	ConfigureResponse(resp, statusCode, payload)
}

// Function to return error response
func RespondError(resp http.ResponseWriter, req *http.Request, statusCode int, errorMessage string) {
	errorResp := struct {
		ErrorMessage string `json: "errorMessage"`
	}{
		ErrorMessage: errorMessage,
	}

	ConfigureResponse(resp, statusCode, errorResp)
}

// Function to configure response
func ConfigureResponse(resp http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("ERROR: unable to marshal data for a response")
		return
	}

	// Configuring the response
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	_, err = resp.Write(data)
	if err != nil {
		message := fmt.Sprintf("ERROR: Unable to write data in response: %v", err)
		fmt.Println(message)
		return
	}
}
