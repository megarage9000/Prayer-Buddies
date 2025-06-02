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

	// 3a. Also send the response of the JSON decode
	prayer, err := config.Database.CreatePrayer(req.Context(), prayerParams)
	if err != nil {
		message := "ERROR, unable to upload prayer to database"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}
	RespondJSON(resp, req, prayer, http.StatusOK)
}

/*
Enables the user to list prayers
TODO: enable optional parameters (i.e. scoping from which users they received from, number of prayers to receive)
*/
func (config *Config) ListPrayers(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab user from header
	userID, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := "ERROR, unable to grab user from header"
		fmt.Println(message)
		RespondError(resp, req, http.StatusUnauthorized, message)
	}

	// TODO: 2. Grab query parameters for prayers from req.URL.Query().Get()

	// 3. Grab the list from user
	listPrayerParams := database.GetReceivedPrayersForUserParams{
		Receiver: userID,
		Limit:    50, // TODO: For now use 50, need to create a constant
	}

	prayers, err := config.Database.GetReceivedPrayersForUser(req.Context(), listPrayerParams)
	if err != nil {
		message := "ERROR, unable to retrieve prayers from database for those receiving"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
	}

	// 3a. Send the prayers to receiver 
	RespondJSON(resp, req, prayers, http.StatusAccepted)
}
