package tracks

import (
	"fmt"
	"molta9a/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListProgressResponse struct {
	UserTracks []db.UserTrack
}

func ListProgress(c *gin.Context) {
	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	join_track := c.DefaultQuery("join_track", "false")

	var userTracks []db.UserTrack

	// Query the database to fetch tracks

	var result *gorm.DB
	if join_track == "true" {
		result = db.DB.Joins("Track").Find(&userTracks).Where("user_id = ?", user.ID)
	} else {
		result = db.DB.Where("user_id = ?", user.ID).Find(&userTracks)
	}
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user tracks"})
		return
	}

	// Prepare the response
	response := ListProgressResponse{
		UserTracks: userTracks,
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
