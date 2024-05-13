package coursesApi

import (
	coursesController "molta9a/server/controllers/courses"
	"molta9a/server/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListResponse struct {
	Courses []db.Course `json:"courses"`
	Total   int         `json:"total"`
}

func List(c *gin.Context) {

	// Query the database to fetch courses
	count, courses, status, err := coursesController.List()
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response
	response := ListResponse{
		Courses: *courses,
		Total:   count,
	}

	// Send the response back to the client
	c.JSON(http.StatusOK, response)
}
