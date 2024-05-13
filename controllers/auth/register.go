package authController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
)

func CreateUser(username, email, password string) (*db.User, int, error) {
	// Check if the username or email already exists
	if UsernameExists(username) {
		return nil, http.StatusConflict, utils.ErrConflictUsername
	}
	if EmailExists(email) {
		return nil, http.StatusConflict, utils.ErrConflictEmail
	}

	// Hash the password before storing it in the database
	hashedPassword, err := HashPassword([]byte(password))
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	user := db.User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}
	return &user, http.StatusCreated, nil
}
