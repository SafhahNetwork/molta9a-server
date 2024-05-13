package userCoursesController

import (
	"errors"
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func Get(userID uint, courseID uint) (int, *db.UserCourse, int, error) {
	var userCourse db.UserCourse

	// Query the database to fetch user course
	var result *gorm.DB
	result = db.DB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&userCourse)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil, http.StatusNotFound, nil
	}

	return int(result.RowsAffected), &userCourse, http.StatusFound, nil
}
