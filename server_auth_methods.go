package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/megarage9000/Prayer-Buddies/internal/auth"
	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

const TOKEN_EXPIRY = time.Hour
const ISSUER = "Prayer Buddies"

/*
Creating a user. This happens on registering users
Response: User with JSON Webtoken (via UserResponse)
On Error: http.StatusInternalServerError
*/
func (config *Config) CreateUser(resp http.ResponseWriter, req *http.Request) {

	// 1. Decoding the response
	decoder := json.NewDecoder(req.Body)
	result := UserCreateRequest{}

	err := decoder.Decode(&result)
	if err != nil {
		message := "ERROR: unable to decode request"
		fmt.Print(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 2. Hashing the password
	hashedPassword, err := auth.HashPassword(result.Password)
	if err != nil {
		message := "ERROR: unable to hash password"
		fmt.Print(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 3. Upload user to database
	userParams := database.RegisterUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		HashedPassword: hashedPassword,
		Email:          result.Email,
	}

	user, err := config.Database.RegisterUser(req.Context(), userParams)
	if err != nil {
		message := "Error: unable to upload user to database"
		fmt.Print(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 4. Creating a JSON WebToken to respond with response
	jsonToken, err := auth.CreateJWT(user.ID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Error: unable to create JSON for user %s", user.ID.String())
		fmt.Print(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 5. Create Response
	responseBody := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jsonToken,
	}

	RespondJSON(resp, req, responseBody, http.StatusOK)
}
