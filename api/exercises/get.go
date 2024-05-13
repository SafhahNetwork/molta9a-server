package exercisesApi

import (
	"log"
	exercisesController "molta9a/server/controllers/exercises"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	Exercise db.Exercise `json:"exercise"`
}

func Get(c *gin.Context) {
	// Extract exercise ID from request parameters
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	// Query the database to fetch the exercise by ID
	found, exercise, status, err := exercisesController.GetByID(uint(exerciseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	if found < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrExerciseNotFound.Error()})
		return
	}

	response := GetResponse{
		Exercise: *exercise,
	}

	c.JSON(http.StatusOK, response)
}
