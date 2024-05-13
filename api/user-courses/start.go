package userCoursesApi

import (
	"log"
	userCoursesController "molta9a/server/controllers/user-courses"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StartResponse struct {
	Created    bool          `json:"created"`
	UserCourse db.UserCourse `json:"user_course"`
}

func Start(c *gin.Context) {
	// Retrieve user information from Gin context
	user := c.MustGet("user").(*db.User)

	// Extract track ID from request parameters
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	ok, userCourse, status, err := userCoursesController.Start(user.ID, uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response
	response := StartResponse{
		Created:    ok,
		UserCourse: *userCourse,
	}

	// Send the user track as a response
	c.JSON(http.StatusCreated, response)
}
