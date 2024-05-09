package auth

import (
	"molta9a/server/db"
	"molta9a/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("NTerneodGaTleltUrLOtAYoNeuTErGenciDeRoPiSteMIe")

// VerifyTok verifies the JWT token from the request header
func VerifyToken(c *gin.Context) (*db.User, error) {
	tokenString, err := c.Cookie("molta9a_biscuit")
	if err != nil || tokenString == "" {
		return nil, utils.ErrMissingAuthorizationToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, utils.ErrInvalidAuthorizationToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.ErrInvalidAuthorizationTokenClaims
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, utils.ErrInvalidAuthorizationTokenClaims
	}

	// Fetch user details from the database using the user ID
	var user db.User
	if err := db.DB.Where("id = ?", uint(userID)).First(&user).Error; err != nil {
		return nil, utils.ErrInvalidAuthorizationTokenClaims
	}

	return &user, nil
}
