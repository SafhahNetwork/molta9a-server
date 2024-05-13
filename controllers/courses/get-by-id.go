package coursesController

import (
	"errors"
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func GetByID(id uint) (int, *db.Course, int, error) {
	// Query the database to fetch the course by ID
	var course db.Course
	result := db.DB.First(&course, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil, http.StatusNotFound, utils.ErrCourseNotFound
		}
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &course, http.StatusFound, nil
}
