package exercisesController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
)

func ListByCourseID(courseID uint) (int, *[]db.Exercise, int, error) {
	// Query to fetch exercises by courseID
	var exercises []db.Exercise
	result := db.DB.Where("course_id = ?", courseID).Order("id ASC").Find(&exercises)

	if result.Error != nil {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &exercises, http.StatusOK, nil
}
