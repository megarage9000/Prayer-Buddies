package main

import "github.com/google/uuid"

type FriendRequest struct {
	ReceivingUser uuid.UUID `json:"receiving_user"`
}

type FriendRequestResponse struct {
	Respondant uuid.UUID `json:"respondant"`
	Status     string    `json: "status"`
}
