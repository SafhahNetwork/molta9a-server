package userCoursesController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func List(userID uint, join_course bool) (int, *[]db.UserCourse, int, error) {
	var userCourses []db.UserCourse

	// Query the database to fetch tracks
	var result *gorm.DB
	if join_course {
		result = db.DB.Joins("Course").Find(&userCourses).Where("user_id = ?", userID)
	} else {
		result = db.DB.Where("user_id = ?", userID).Find(&userCourses)
	}

	if result.Error != nil {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &userCourses, http.StatusOK, nil
}
