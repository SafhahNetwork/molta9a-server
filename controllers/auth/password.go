package authController

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	return string(hash), err
}

func CheckPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
