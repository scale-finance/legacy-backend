package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// encrypts password with all appropriate settings and conversions for simple use in the
// main authentication file
func EncryptPassword(password string) string {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encrypted)
}

// The purpose of this function is to simplify the code in the authentication file. It works
// the same as bcrypt's library function, but with preset settings already integrated in the
// function call itself.
func HashMatch(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	} else {
		return true
	}
}