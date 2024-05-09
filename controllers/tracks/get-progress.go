package tracks

import (
	"fmt"
	"log"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetProgressResponse struct {
	UserModules []db.UserModule
}

func GetProgress(c *gin.Context) {
	// Extract track ID from request parameters
	trackIDStr := c.Param("id")
	trackID, err := strconv.ParseUint(trackIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)
	var userModules []db.UserModule

	// Query the database to fetch tracks
	result := db.DB.Where("user_id = ? AND track_id = ?", user.ID, trackID).Find(&userModules)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user modules"})
		return
	}

	// Prepare the response
	response := GetProgressResponse{
		UserModules: userModules,
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
