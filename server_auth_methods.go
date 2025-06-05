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
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 2. Hashing the password
	hashedPassword, err := auth.HashPassword(result.Password)
	if err != nil {
		message := "ERROR: unable to hash password"
		fmt.Println(message)
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
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 4. Creating a JSON WebToken to respond with response
	jsonToken, err := auth.CreateJWT(user.ID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Error: unable to create JSON for user %s", user.ID.String())
		fmt.Println(message)
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

/*
Logging in a user. This happens on registering users
Response: User with JSON Webtoken (via UserResponse)
On Error: http.StatusInternalServerError, http.StatusUnauthorized
*/
func (config *Config) LoginUser(resp http.ResponseWriter, req *http.Request) {

	// 1. Decoding the response
	decoder := json.NewDecoder(req.Body)
	result := UserCreateRequest{}

	err := decoder.Decode(&result)
	if err != nil {
		message := "ERROR: unable to decode request"
		fmt.Println(message)
		RespondError(resp, req, http.StatusInternalServerError, message)
		return
	}

	// 2. Check if the user exists
	user, err := config.Database.GetUserByEmail(req.Context(), result.Email)
	if err != nil {
		message := "ERROR: unable to get user from database"
		fmt.Println(message)
		RespondError(resp, req, http.StatusUnauthorized, message)
		return
	}

	// 3. Hash the received password, and check if it matches the one on database
	err = auth.CheckPasswordWithHash(result.Password, user.HashedPassword)
	if err != nil {
		message := fmt.Sprintf("ERROR: invalid password: %s", err)
		fmt.Println(message)
		RespondError(resp, req, http.StatusUnauthorized, message)
		return
	}

	// 4. Return a JSON Webtoken on successful login
	jsonToken, err := auth.CreateJWT(user.ID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Error: unable to create JSON for user %s", user.ID.String())
		fmt.Println(message)
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
