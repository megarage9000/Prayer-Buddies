package main

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateRequest struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

type UserSetUsernameRequest struct {
	Username string `json: "username"`
}

type UserResponse struct {
	ID        uuid.UUID `json: "id"`
	CreatedAt time.Time `json: "created_at"`
	UpdatedAt time.Time `json: "updated_at"`
	Email     string    `json: "email"`
	Token     string    `json: "json_token"`
}
