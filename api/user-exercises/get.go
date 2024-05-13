package userExercisesApi

import (
	"log"
	"molta9a/server/controllers/user-exercises"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	UserExercise db.UserExercise `json:"user_exercise"`
}

func Get(c *gin.Context) {
	// Extract exercise ID from request parameters
	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	// Query the database to fetch the user exercise by ID
	found, userExercise, status, err := userExercisesController.Get(user.ID, uint(exerciseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	if found < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrDidntStartExercise})
		return
	}

	// Prepare the response
	response := GetResponse{
		UserExercise: *userExercise,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
