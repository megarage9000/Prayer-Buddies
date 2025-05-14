package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

// Testing password hashing
func TestHashPassword(t *testing.T) {

	password := "hello123"

	passwordBytes, err := HashPassword(password)

	if err != nil {
		t.Fatalf("FAILED: Error occured in trying to hash password: %v", err)
	} else if passwordBytes == password {
		t.Fatalf("FAILED: Resulting password did not result in a hash")
	}
}

// Test JWT Creation
func TestCreateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "my dirty little secret"
	issuer := "tester"
	expiry := time.Second * 2

	_, err := CreateJWT(userID, secret, issuer, expiry)

	if err != nil {
		t.Fatalf("FAILED: Error occured in creating JWT: %v", err)
	}
}

// Test if validation works on valid expiry tokens
func TestValidateJWTNonExpiry(t *testing.T) {
	userID := uuid.New()
	secret := "my dirty little secret"
	issuer := "tester"
	expiry := time.Second * 2

	token, err := CreateJWT(userID, secret, issuer, expiry)

	if err != nil {
		t.Fatalf("FAILED: Error occured in creating JWT: %v", err)
	}

	// wait one second
	time.Sleep(time.Second)

	resultID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("FAILED: Error occured in validating a valid JWT: %v", err)
	} else if resultID != userID {
		t.Fatalf("FAILED: Expected result ID to be %v, got %v", userID, resultID)
	}
}

// Test if validation fails on expired tokens
func TestValidateJWTExpiry(t *testing.T) {
	userID := uuid.New()
	secret := "my dirty little secret"
	issuer := "tester"
	expiry := time.Second * 2

	token, err := CreateJWT(userID, secret, issuer, expiry)

	if err != nil {
		t.Fatalf("FAILED: Error occured in creating JWT: %v", err)
	}

	// wait 3 seconds
	time.Sleep(time.Second * 3)

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("FAILED: Validation did not fail on expired token")
	}
}
