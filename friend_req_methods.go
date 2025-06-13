package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

func (config *Config) SendFriendRequest(resp http.ResponseWriter, req *http.Request) {

	//1. Grab the user from the JWT token
	user_id, err := GrabUserIDFromHeader(req.Header, *config)

	if err != nil {
		message := "Unable to grab user from from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Retrieve the friend id
	decoder := json.NewDecoder(req.Body)
	result := FriendRequest{}

	err = decoder.Decode(&result)

	if err != nil {
		message := "Unable to decode friend request"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 3. Upload the request onto the database
	friendRequestParams := database.CreateFriendRequestParams{
		UserID:    user_id,
		FriendID:  result.ReceivingUser,
		CreatedAt: time.Now(),
	}

	friendReqId, err := config.Database.CreateFriendRequest(req.Context(), friendRequestParams)
	if err != nil {
		message := "Unable to upload friend request to database"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	RespondJSON(resp, req, friendReqId, http.StatusOK)
}
