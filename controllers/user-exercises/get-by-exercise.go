package userExercisesController

import (
	"errors"
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func Get(userID uint, exerciseID uint) (int, *db.UserExercise, int, error) {
	var userExercise db.UserExercise

	// Query the database to fetch user exercise
	var result *gorm.DB
	result = db.DB.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).First(&userExercise)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil, http.StatusNotFound, nil
	}

	return int(result.RowsAffected), &userExercise, http.StatusFound, nil
}
