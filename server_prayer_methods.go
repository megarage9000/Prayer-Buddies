package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

/*
Sends a prayer request, user must be logged in with a JWT
*/
func (config *Config) SendPrayerRequest(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab user ID from JWT token
	sendingUserID, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := fmt.Sprintf("ERROR, unable to grab user ID from header", err)
		fmt.Println(message)
		RespondError(resp, req, http.StatusUnauthorized, message)
		return
	}

	// 2. Grab prayer request on body
	decoder := json.NewDecoder(req.Body)
	result := PrayerRequest{}

	err = decoder.Decode(&result)
	if err != nil {
		message := "ERROR, unable to decode error"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 3. Send the prayer request on the database
	receiver, err := uuid.Parse(result.Receiver)
	if err != nil {
		message := "ERROR, unable to parse receiver UUID"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	prayerParams := database.CreatePrayerParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sender:    sendingUserID,
		Receiver:  receiver,
		Prayer:    result.Prayer,
	}

	prayer, err := config.Database.CreatePrayer(req.Context(), prayerParams)
	if err != nil {
		message := "ERROR, unable to upload prayer to database"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}
	RespondJSON(resp, req, prayer, http.StatusOK)
}
