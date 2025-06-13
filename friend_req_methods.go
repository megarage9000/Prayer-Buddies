package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

const (
	ACCEPTED = "accepted"
	DENIED   = "denied"
	PENDING  = "pending"
)

/*
	Sends a friend request on the server, which will be stored on the database. The status will be stored as "pending"
*/

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

/*
Updates the friend request accordingly
- If accepted, set status to accepted and add another row for the friend
- If denied, remove the request
*/
func (config *Config) UpdateFriendRequest(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab the user
	user, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := "Unable to get user from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Decode response
	decoder := json.NewDecoder(req.Body)
	result := FriendRequestResponse{}
	err = decoder.Decode(&result)
	if err != nil {
		message := "Unable to decode friend request response"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 3. Check response and updating request
	status := result.Status

	updateFriendRequestArgs := database.UpdateFriendRequestParams{
		UserID:   result.Respondant,
		FriendID: user,
		Status:   status,
	}

	err = config.Database.UpdateFriendRequest(req.Context(), updateFriendRequestArgs)
	if err != nil {
		message := "Could not update friend request status"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 4. Update the result accordingly
	switch status {
	case ACCEPTED:
		{
			acceptFriendReq := database.AcceptFriendRequestParams{
				UserID:    user,
				FriendID:  result.Respondant,
				CreatedAt: time.Now(),
			}

			res, err := config.Database.AcceptFriendRequest(req.Context(), acceptFriendReq)
			if err != nil {
				message := "Unable to add new row to represent new friend"
				LogError(message, err, resp, req, http.StatusInternalServerError)
				return
			}

			RespondJSON(resp, req, res, http.StatusAccepted)
		}
	case DENIED:
		{
			denyFriendReq := database.DenyFriendRequestParams{
				UserID:   result.Respondant,
				FriendID: user,
			}

			err = config.Database.DenyFriendRequest(req.Context(), denyFriendReq)
			if err != nil {
				message := "Unable to remove pending request to represent denial"
				LogError(message, err, resp, req, http.StatusInternalServerError)
				return
			}
		}
	}

}
