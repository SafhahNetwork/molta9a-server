package selfApi

import (
	selfController "molta9a/server/controllers/self"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	User db.User `json:"user"`
	Rank int     `json:"rank"`
}

func Get(c *gin.Context) {
	// Retrieve user information from Gin context
	cuser := c.MustGet("user").(*db.User)

	ok, user, rank, status, err := selfController.Get(cuser.ID)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrUserNotFound.Error()})
		return
	}

	response := GetResponse{
		User: *user,
		Rank: rank,
	}

	c.JSON(http.StatusOK, response)
}
