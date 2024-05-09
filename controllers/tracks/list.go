package tracks

import (
	"fmt"
	"molta9a/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	Tracks []db.Track `json:"tracks"`
	Total  int64      `json:"total"`
}

func List(c *gin.Context) {
	var tracks []db.Track

	// Query the database to fetch tracks
	result := db.DB.Find(&tracks)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tracks"})
		return
	}

	// Prepare the response
	response := ListResponse{
		Tracks: tracks,
		Total:  result.RowsAffected,
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
