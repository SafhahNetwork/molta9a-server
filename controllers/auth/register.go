package auth

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Username string `form:"username" binding:"required,min=4,max=64"`
	Password string `form:"password" binding:"required,min=6,max=2048"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		c.Error(err)
		return
	}

	// Check if the username or email already exists
	if UsernameExists(req.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": utils.ErrConflictUsername.Error()})
		return
	}
	if EmailExists(req.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": utils.ErrConflictEmail.Error()})
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := HashPassword([]byte(req.Password))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	user := db.User{
		Email:          req.Email,
		Username:       req.Username,
		Password:       hashedPassword,
		TotalXp:        0,
		TotalExercises: 0,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	// Generate JWT and send it to the user
	tokenString, err := generateToken(user.ID, time.Hour*STAY_LOGGED_HOURS)
	if err != nil {
		log.Println(err)
		tokenString = ""
	}

	// Set the JWT token as an HttpOnly cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("molta9a_biscuit", tokenString, 60*60*STAY_LOGGED_HOURS, "/", "localhost", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": utils.InfoUserCreated,
		"id":      user.ID,
		"expires": time.Now().Add(time.Hour * 24),
	})
}
