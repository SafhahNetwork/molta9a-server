package auth

import (
	"molta9a/server/db"
)

func UsernameExists(username string) bool {
	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return true
}

func EmailExists(email string) bool {
	var user db.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}
