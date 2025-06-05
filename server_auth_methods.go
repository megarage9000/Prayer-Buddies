package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/megarage9000/Prayer-Buddies/internal/auth"
	"github.com/megarage9000/Prayer-Buddies/internal/database"
)

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
		message := "Unable to decode request"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 2. Hashing the password
	hashedPassword, err := auth.HashPassword(result.Password)
	if err != nil {
		message := "Unable to hash password"
		LogError(message, err, resp, req, http.StatusInternalServerError)
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
		message := "Unable to upload user to database"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 4. Creating a JSON WebToken to respond with response
	jsonToken, err := auth.CreateJWT(user.ID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Unable to create JSON for user %s", user.ID.String())
		LogError(message, err, resp, req, http.StatusInternalServerError)
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
		message := "Unable to decode request"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 2. Check if the user exists
	user, err := config.Database.GetUserByEmail(req.Context(), result.Email)
	if err != nil {
		message := "Unable to get user from database"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 3. Hash the received password, and check if it matches the one on database
	err = auth.CheckPasswordWithHash(result.Password, user.HashedPassword)
	if err != nil {
		message := "Invalid password"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 4. Return a JSON Webtoken on successful login
	jsonToken, err := auth.CreateJWT(user.ID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Unable to create JSON for user %s", user.ID.String())
		LogError(message, err, resp, req, http.StatusInternalServerError)
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
For a given refresh token, refreshes to return a new JWT
*/
func (config *Config) RefreshToken(resp http.ResponseWriter, req *http.Request) {
	// 1. Grab refresh token from header
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		message := "Unable to get refresh token from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Check if the refresh token exists
	result, err := config.Database.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		message := "Unable to get refresh token from database, might be expired or revoked"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
		// TODO: Do we make a new refresh token for the user if it fails?
	}

	// 3. Create a new JWT for user
	userID := result.UserID
	jsonToken, err := auth.CreateJWT(userID, config.Secret, ISSUER, TOKEN_EXPIRY)
	if err != nil {
		message := fmt.Sprintf("Unable to create JSON for user %s", userID.String())
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 4. Return JWT token
	payload := struct {
		token string
	} {
		token: jsonToken,
	}

	RespondJSON(resp, req, payload, http.StatusOK)
}

/*
Revokes the current refresh token, for a given token
*/
func (config *Config) RevokeToken(resp http.ResponseWriter, req *http.Request) {

	// 1. Grab refresh token from header
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		message := "Unable to get refresh token from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Update the database to revoke the refresh token
	revokeTokenParams := database.RevokeTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = config.Database.RevokeToken(req.Context(), revokeTokenParams)
	if err != nil {
		message := "Unable to revoke token from header"
		LogError(message, err, resp, req, http.StatusInternalServerError)
	}

	// 3. Return a no content header
	resp.WriteHeader(http.StatusNoContent)
}

/*
Helper Function that generates a refresh token for the user
*/
func (config *Config) CreateRefreshTokenForUser(userID uuid.UUID, ctx context.Context) (string, error) {

	// 1. Get Refresh Token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}

	// 2. Upload Refresh Token to database for user
	createRefreshToken := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24), // For now refresh tokens last a day!
		UserID:    userID,
	}

	_, err = config.Database.CreateRefreshToken(ctx, createRefreshToken)
	if err != nil {
		return "", err
	}

	// 3. Return Refresh Token
	return refreshToken, nil
}

/*
Set username
*/

func (config *Config) SetUsername(resp http.ResponseWriter, req *http.Request) {

	// 1. Get User from server
	user, err := GrabUserIDFromHeader(req.Header, *config)
	if err != nil {
		message := "Unable to grab user from header"
		LogError(message, err, resp, req, http.StatusUnauthorized)
		return
	}

	// 2. Get username data
	decoder := json.NewDecoder(req.Body)
	result := UserSetUsernameRequest{}

	err = decoder.Decode(&result)
	if err != nil {
		message := "Unable to decode set username response"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	// 3. Update user with data
	usernameParams := database.SetUsernameParams{
		ID:       user,
		Username: result.Username,
	}

	err = config.Database.SetUsername(req.Context(), usernameParams)
	if err != nil {
		message := "Unable to set username for user"
		LogError(message, err, resp, req, http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusAccepted)
}
