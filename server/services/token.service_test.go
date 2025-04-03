package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAccessToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    userID := 123
    token, err := service.GenerateAccessToken(userID)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if token == nil || *token == "" {
        t.Fatalf("expected a valid token, got nil or empty string")
    }
}

func TestValidateToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    userID := 123
    token, err := service.GenerateAccessToken(userID)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    // Validate the generated token
    validatedUserID, err := service.ValidateToken(token)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    if validatedUserID == nil || *validatedUserID != userID {
        t.Fatalf("expected userID %d, got %v", userID, validatedUserID)
    }
}

func TestValidateToken_InvalidToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    invalidToken := "invalid.token.value"
    _, err := service.ValidateToken(&invalidToken)
    if err == nil {
        t.Fatalf("expected an error for invalid token, got nil")
    }

    expectedError := "JWT could not be validated"
    if err.Error() != expectedError {
        t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
    }
}

func TestValidateToken_ExpiredToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    // Create an expired token
    expiredTime := time.Now().Add(-1 * time.Hour).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
        "iss":  claimIssuer,
        "sub":  "api_access_token",
        "exp":  expiredTime,
        "user": 123,
    })
    signedToken, err := token.SignedString([]byte(jwtSecretKey))
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    // Validate the expired token
    _, err = service.ValidateToken(&signedToken)
    if err == nil {
        t.Fatalf("expected an error for expired token, got nil")
    }

    expectedError := "JWT could not be validated"
    if err.Error() != expectedError {
        t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
    }
}

func TestValidateToken_MalformedToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    // Malformed token (missing parts)
    malformedToken := "malformed.token"
    _, err := service.ValidateToken(&malformedToken)
    if err == nil {
        t.Fatalf("expected an error for malformed token, got nil")
    }

    expectedError := "JWT could not be validated"
    if err.Error() != expectedError {
        t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
    }
}

func TestValidateToken_EmptyToken(t *testing.T) {
    jwtSecretKey := "testSecretKey"
    service := NewTokenService(jwtSecretKey)

    // Empty token
    emptyToken := ""
    _, err := service.ValidateToken(&emptyToken)
    if err == nil {
        t.Fatalf("expected an error for empty token, got nil")
    }

    expectedError := "JWT could not be validated"
    if err.Error() != expectedError {
        t.Fatalf("expected error '%s', got '%s'", expectedError, err.Error())
    }
}