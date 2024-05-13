package authApi

import (
	"log"
	authController "molta9a/server/controllers/auth"
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

	user, status, err := authController.CreateUser(req.Username, req.Email, req.Password)
	if err != nil || user == nil {
		log.Println(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT and send it to the user
	tokenString, err := authController.GenerateToken(user.ID)
	if err != nil {
		log.Println(err)
		tokenString = ""
	}

	// Set the JWT token as an HttpOnly cookie
	authController.SetGinCookie(c, tokenString)

	c.JSON(http.StatusCreated, gin.H{
		"message": utils.InfoUserCreated,
		"id":      user.ID,
		"expires": time.Now().Add(time.Hour * authController.STAY_LOGGED_HOURS),
	})
}
