package userCoursesApi

import (
	userCoursesController "molta9a/server/controllers/user-courses"
	"molta9a/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	UserCourses []db.UserCourse `json:"user_courses"`
}

func List(c *gin.Context) {
	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	// Query the database to fetch the user course by ID
	_, userCourses, status, err := userCoursesController.List(user.ID, true)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response
	response := ListResponse{
		UserCourses: *userCourses,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
