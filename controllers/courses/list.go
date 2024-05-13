package coursesController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
)

func List() (int, *[]db.Course, int, error) {
	// Query to fetch courses
	var courses []db.Course
	result := db.DB.Find(&courses).Order("id ASC")

	if result.Error != nil {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &courses, http.StatusOK, nil
}
