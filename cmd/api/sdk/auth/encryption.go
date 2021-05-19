package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) string {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encrypted)
}

func HashMatch(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}