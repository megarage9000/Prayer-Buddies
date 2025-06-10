package main

import (
	"encoding/json"
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
		message := "Unable to grab user ID from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Grab prayer request on body
	decoder := json.NewDecoder(req.Body)
	result := PrayerRequest{}

	err = decoder.Decode(&result)
	if err != nil {
		message := "Unable to decode error"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 3. Send the prayer request on the database
	receiver, err := uuid.Parse(result.Receiver)
	if err != nil {
		message := "Unable to parse receiver UUID"
		LogError(message, err, resp, req, http.StatusInternalServerError)
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
		message := "Unable to upload prayer to database"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}
	RespondJSON(resp, req, prayer, http.StatusOK)
}

/*
Enables the user to list prayer requests received
TODO: enable optional parameters (i.e. scoping from which users they received from, number of prayers to receive)
*/
func (config *Config) ListReceivedPrayerRequests(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab user from header
	userID, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := "Unable to grab user from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// TODO: 2. Grab query parameters for prayers from req.URL.Query().Get()

	// 3. Grab the list from user
	listPrayerParams := database.GetReceivedPrayersForUserParams{
		Receiver: userID,
		Limit:    50, // TODO: For now use 50, need to create a constant
	}

	prayers, err := config.Database.GetReceivedPrayersForUser(req.Context(), listPrayerParams)
	if err != nil {
		message := "Unable to retrieve prayers from database for those receiving"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 3a. Send the prayers to receiver
	RespondJSON(resp, req, prayers, http.StatusAccepted)
}

/*
Enables the user to list prayer requests sent
TODO: enable optional parameters (i.e. scoping from which users they sent to, number of prayers sent)
*/
func (config *Config) ListSentPrayerRequests(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab user from header
	userID, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := "Unable to grab user from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// TODO: 2. Grab query parameters for prayers from req.URL.Query().Get()

	// 3. Grab list from user
	prayersSentParams := database.GetSentPrayersFromUserParams{
		Sender: userID,
		Limit:  50, // TODO: For now use 50, need to create a constant
	}

	prayers, err := config.Database.GetSentPrayersFromUser(req.Context(), prayersSentParams)
	if err != nil {
		message := "Unable to retrieve prayers from database for those sending"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 3a. Sending the sent prayer requests to user
	RespondJSON(resp, req, prayers, http.StatusOK)
}
