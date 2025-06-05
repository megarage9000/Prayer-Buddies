package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/megarage9000/Prayer-Buddies/internal/auth"
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
	resp.Write(data)
}

// Function to grab a userID from http header
func GrabUserIDFromHeader(header http.Header, config Config) (uuid.UUID, error) {

	token, err := auth.GetBearerToken(header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ERROR: unable to get bearer token")
	}

	user, err := auth.ValidateJWT(token, config.Secret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ERROR: unable to validate jwt")
	}

	return user, nil
}

func LogError(message string, err error, resp http.ResponseWriter, req *http.Request, statusCode int) {
	logLine := fmt.Sprintf("ERROR: %s\n - error details: %v\n", message, err)
	fmt.Println(logLine)
	RespondError(resp, req, statusCode, message)
}
