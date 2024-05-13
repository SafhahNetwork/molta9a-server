package userCoursesController

import (
	"log"
	coursesController "molta9a/server/controllers/courses"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
)

func Start(userID uint, courseID uint) (bool, *db.UserCourse, int, error) {
	_, course, status, err := coursesController.GetByID(uint(courseID))
	if err != nil {
		return false, nil, status, err
	}

	found, oldUserCourse, status, err := Get(userID, courseID)
	if err != nil {
		return false, nil, status, err
	}
	if found > 0 {
		return false, oldUserCourse, http.StatusConflict, utils.ErrAlreadyStartedCourse
	}

	userCourse := db.UserCourse{
		CourseID: course.ID,
		UserID:   userID,
		Status:   "active",
		Progress: 0,
	}
	result := db.DB.Create(&userCourse)
	if result.Error != nil {
		log.Println(err)
		return false, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return true, &userCourse, http.StatusCreated, nil
}
