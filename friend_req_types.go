package main

import "github.com/google/uuid"

type FriendRequest struct {
	ReceivingUser uuid.UUID `json:"receiving_user"`
}
