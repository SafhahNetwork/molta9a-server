package authApi

import (
	"log"
	authController "molta9a/server/controllers/auth"
	"molta9a/server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
	user, err := authController.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrBadLogin.Error()})
		return
	}

	tokenString, err := authController.GenerateToken(user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrInternalServerError.Error()})
		return
	}

	// Set the JWT token as an HttpOnly cookie
	authController.SetGinCookie(c, tokenString)

	c.JSON(http.StatusOK, gin.H{
		"message": utils.InfoSuccess,
		"id":      user.ID,
		"expires": time.Now().Add(time.Hour * 24),
	})
}
