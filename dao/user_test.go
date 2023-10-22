package dao

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestUserBeforeSave(t *testing.T) {
	user := &User{
		Password: "password123",
	}

	err := user.BeforeSave(nil)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if user.Password == "password123" {
		t.Errorf("Password should have been hashed, but it wasn't")
	}
}

func TestUserCheckPassword(t *testing.T) {
	stringPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(stringPassword), bcrypt.DefaultCost)
	user := &User{
		Password: string(hashedPassword),
	}

	err := user.CheckPassword(stringPassword)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	err = user.CheckPassword("wrongpassword")
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}
