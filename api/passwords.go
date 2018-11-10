package main

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a string and returns a hash and possible error
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash takes a password and hash string, and returns if the hash came from password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
