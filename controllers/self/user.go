package selfController

import (
	"errors"
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"

	"gorm.io/gorm"
)

func Get(userID uint) (bool, *db.User, int, int, error) {
	var fuser db.User

	// Query the database to fetch tracks
	result := db.DB.
		Select("id", "username", "email", "avatar", "total_xp", "total_exercises", "total_lessons", "total_tracks").
		First(&fuser, userID)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		return false, nil, 0, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	var rank int = 0
	result = db.DB.Raw(`
    SELECT ranked.rank from ( select id, RANK() OVER (ORDER BY total_xp DESC) from users ) AS ranked WHERE id = ?
    `, userID).Scan(&rank)
	if result.Error != nil {
		log.Println(result.Error)
	}

	return true, &fuser, rank, http.StatusFound, nil
}
