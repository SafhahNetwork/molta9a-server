package authController

import (
	"molta9a/server/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	// Create JWT token claims
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(STAY_LOGGED_HOURS * time.Hour).Unix()

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func AuthenticateUser(username, password string) (*db.User, error) {
	var user db.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err // User not found
	}

	// Compare hashed password with input password
	if err := CheckPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err // Password does not match
	}

	return &user, nil // Authentication successful
}
