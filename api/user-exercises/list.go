package userExercisesApi

import (
	"log"
	userExercisesController "molta9a/server/controllers/user-exercises"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	UserExercises []db.UserExercise `json:"user_Exercises"`
}

func List(c *gin.Context) {
	// Extract course ID from request parameters
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	// Query the database to fetch the user course by ID
	_, userExercises, status, err := userExercisesController.List(user.ID, uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response
	response := ListResponse{
		UserExercises: *userExercises,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
