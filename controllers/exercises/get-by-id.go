package exercisesController

import (
	"errors"
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func GetByID(id uint) (int, *db.Exercise, int, error) {
	// Query the database to fetch the Exercise by ID
	var exercise db.Exercise
	result := db.DB.First(&exercise, id)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &exercise, http.StatusFound, nil
}
