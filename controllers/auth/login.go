package auth

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const STAY_LOGGED_HOURS = 24 * 7

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
	cookie, err := c.Cookie("molta9a_biscuit")
	log.Printf("Cookie value: %s \n", cookie)

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	// Authenticate user
	user, err := authenticateUser(req.Username, req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrBadLogin.Error()})
		return
	}

	tokenString, err := generateToken(user.ID, time.Hour*STAY_LOGGED_HOURS)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	// Set the JWT token as an HttpOnly cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("molta9a_biscuit", tokenString, 60*60*STAY_LOGGED_HOURS, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.InfoSuccess,
		"id":      user.ID,
		"expires": time.Now().Add(time.Hour * 24),
	})
}

func generateToken(userID uint, timeout time.Duration) (string, error) {
	// Create JWT token claims
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(timeout).Unix() // Token expires in 24 hours

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func authenticateUser(username, password string) (*db.User, error) {
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
