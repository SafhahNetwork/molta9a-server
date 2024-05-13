package middleware

import (
	"log"
	"net/http"

	authController "molta9a/server/controllers/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware function to verify JWT token and set user information in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authController.VerifyToken(c)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user information in Gin context
		c.Set("user", user)

		c.Next()
	}
}
