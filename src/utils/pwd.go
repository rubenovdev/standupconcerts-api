package utils

import (
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratePassword() (string, error) {
	password, err := password.Generate(10, 5, 2, false, false)
	return password, err
}
