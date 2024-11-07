package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (error, string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return err, string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
