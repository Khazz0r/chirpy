package auth

import (
	"testing"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPasswordAndPlainText(t *testing.T) {
	testPassword := "hello"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Error happened with generating hash from password: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err != nil {
		t.Errorf("Test password does not match the hashed password: %v", err)
	}
}

func TestTwoHashedPasswords(t *testing.T) {
	testPassword := "Goodbye"

	// hash both passwords
	hashedPassword1, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Error happened with generating first hash from password: %v", err)
	}
	hashedPassword2, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("Error happened with generating second password from password: %v", err)
	}

	// compare hash and password to ensure they match
	if err := bcrypt.CompareHashAndPassword(hashedPassword1, []byte(testPassword)); err != nil {
		t.Errorf("Password did not match hashedPassword1: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(hashedPassword2, []byte(testPassword)); err != nil {
		t.Errorf("Password did not match hashedPassword2: %v", err)
	}
}
