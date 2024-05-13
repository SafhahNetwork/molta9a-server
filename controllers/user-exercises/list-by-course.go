package userExercisesController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func List(userID uint, courseID uint) (int, *[]db.UserExercise, int, error) {
	var userExercises []db.UserExercise

	// Query the database to fetch tracks
	var result *gorm.DB
	result = db.DB.Where("user_id = ? AND course_id = ?", userID, courseID).Find(&userExercises)

	if result.Error != nil {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &userExercises, http.StatusOK, nil
}
