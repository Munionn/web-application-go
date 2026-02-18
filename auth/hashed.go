package auth

import (
	"bytes"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytesHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytesHash), nil
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}
