package tracks

import (
	"log"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StartResponse struct {
}

func Start(c *gin.Context) {
	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	// Extract track ID from request parameters
	trackIDStr := c.Param("id")
	trackID, err := strconv.ParseUint(trackIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	// Query the database to fetch the track by ID
	var track db.Track
	result := db.DB.First(&track, trackID)
	if result.Error != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	userTrack := db.UserTrack{
		TrackID:  track.ID,
		UserID:   user.ID,
		Status:   "active",
		Progress: 0,
	}
	result = db.DB.Create(&userTrack)
	if result.Error != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Send the user track as a response
	c.JSON(http.StatusOK, userTrack)
}
