package tracks

import (
	"log"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	Track   db.Track    `json:"track"`
	Quests  []db.Quest  `json:"quests"`
	Modules []db.Module `json:"modules"`
}

func Get(c *gin.Context) {
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

	// Query to fetch quests by trackID
	var quests []db.Quest
	result = db.DB.Where("track_id = ?", trackID).Order("id ASC").Find(&quests)
	if result.Error != nil {
		log.Println(err)
	}

	// Query to fetch modules by trackID
	var modules []db.Module
	result = db.DB.Where("track_id = ?", trackID).Order("index ASC").Find(&modules)
	if result.Error != nil {
		log.Println(err)
	}

	// Prepare the response
	response := GetResponse{
		Track:   track,
		Quests:  quests,
		Modules: modules,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
