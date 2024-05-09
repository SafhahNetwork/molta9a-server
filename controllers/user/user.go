package user

import (
	"log"
	"molta9a/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	User db.User
	Rank int
}

func PrivateGet(c *gin.Context) {
	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	var fuser db.User

	// Query the database to fetch tracks
	result := db.DB.
		Select("id", "username", "email", "avatar", "total_xp", "total_exercises", "total_lessons", "total_tracks").
		First(&fuser, user.ID)
	if result.Error != nil {
		log.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	var rank int
	result = db.DB.Raw(`
    SELECT ranked.rank from ( select id, RANK() OVER (ORDER BY total_xp DESC) from users ) AS ranked WHERE id = ?
    `, user.ID).Scan(&rank)
	if result.Error != nil {
		log.Println(result.Error)
	}

	// Prepare the response
	response := GetResponse{
		User: fuser,
		Rank: rank,
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
