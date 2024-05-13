package coursesApi

import (
	"log"
	coursesController "molta9a/server/controllers/courses"
	exercisesController "molta9a/server/controllers/exercises"
	questsController "molta9a/server/controllers/quests"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetResponse struct {
	Course    db.Course     `json:"course"`
	Quests    []db.Quest    `json:"quests"`
	Exercises []db.Exercise `json:"exercises"`
}

func Get(c *gin.Context) {
	// Extract course ID from request parameters
	courseIDStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Query the database to fetch the course by ID
	_, course, status, err := coursesController.GetByID(uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Query to fetch quests by course
	_, quests, status, err := questsController.ListByCourseID(uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Query to fetch exercises by course
	_, exercises, status, err := exercisesController.ListByCourseID(uint(courseID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response
	response := GetResponse{
		Course:    *course,
		Quests:    *quests,
		Exercises: *exercises,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
