package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --- Hashing ---

// Hash passwords using Bcrypt
func HashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("ERROR: Unable to generate password")
	}
	return string(passwordBytes), nil
}

func CheckPasswordWithHash(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// --- JSON Web Tokens ---

// Creating a JWT
func CreateJWT(userID uuid.UUID, secret string, issuer string, expiry time.Duration) (string, error) {

	// Creating a registered claims
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		Subject:   userID.String(),
	}

	// Creating a token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Returning token with secret
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (uuid.UUID, error) {
	type claimsStruct struct {
		jwt.RegisteredClaims
	}

	var userID uuid.UUID

	// 1. Parsing the Token to ensure it is signed correctly
	extractedToken, err := jwt.ParseWithClaims(tokenString, &claimsStruct{}, func(token *jwt.Token) (interface{}, error) {
		// Type assserting to ensure that the token's Method is of SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ERROR: Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return userID, fmt.Errorf("ERROR: Unable to parse token: %v", err)
	}

	// 2. Checking if we can get Subject Field
	userIDStr, err := extractedToken.Claims.GetSubject()
	if err != nil {
		return userID, fmt.Errorf("ERROR: Unable to get Subject field from token: %v", err)
	}

	// 3. Checking if we can parse the Subject Field into a UUID
	userID, err = uuid.Parse(userIDStr)
	if err != nil {
		return userID, fmt.Errorf("ERROR: Unable to parse %s to uuid.UUID: %v", userIDStr, err)
	}

	return userID, nil
}
