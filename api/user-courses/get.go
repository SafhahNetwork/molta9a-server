package userCoursesApi

import (
	"log"
	userCoursesController "molta9a/server/controllers/user-courses"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	UserCourse db.UserCourse `json:"user_course"`
}

func Get(c *gin.Context) {
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
	found, userCourse, status, err := userCoursesController.Get(user.ID, uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	if found < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrDidntStartCourse})
		return
	}

	// Prepare the response
	response := GetResponse{
		UserCourse: *userCourse,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
