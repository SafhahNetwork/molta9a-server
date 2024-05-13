package questsController

import (
	"log"
	"molta9a/server/db"
	"molta9a/server/utils"
	"net/http"
)

func ListByCourseID(courseID uint) (int, *[]db.Quest, int, error) {
	// Query to fetch quests by trackID
	var quests []db.Quest
	result := db.DB.Where("course_id = ?", courseID).Order("id ASC").Find(&quests)

	if result.Error != nil {
		log.Println(result.Error)
		return 0, nil, http.StatusInternalServerError, utils.ErrInternalServerError
	}

	return int(result.RowsAffected), &quests, http.StatusOK, nil
}
