package auth

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	auth := AuthInit()
	userId := "123"
	expiresIn, tokenString, err := auth.GenerateJWT(userId)

	if err != nil {
		t.Errorf("Error while generating JWT: %v", err)
	}

	if expiresIn <= 0 {
		t.Error("ExpiresIn should be greater than 0")
	}

	if tokenString == "" {
		t.Error("TokenString should not be empty")
	}
}

func TestValidateTokenValid(t *testing.T) {
	auth := AuthInit()
	userId := "123"
	_, tokenString, err := auth.GenerateJWT(userId)
	if err != nil {
		t.Fatalf("Error while generating JWT: %v", err)
	}

	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		t.Errorf("Error while validating token: %v", err)
	}

	if claims.UserID != userId {
		t.Errorf("Expected UserID: %s, got: %s", userId, claims.UserID)
	}

	if claims.ExpiresAt <= time.Now().Unix() {
		t.Error("Token should not be expired")
	}
}

func TestValidateTokenInvalid(t *testing.T) {
	auth := AuthInit()
	tokenString := "invalid-token-string"
	_, err := auth.ValidateToken(tokenString)
	if err == nil {
		t.Error("Token validation should fail for invalid token string")
	}
}

func TestValidateTokenEmpty(t *testing.T) {
	auth := AuthInit()
	tokenString := ""
	_, err := auth.ValidateToken(tokenString)
	if err == nil {
		t.Error("Token validation should fail for empty token string")
	}
}
