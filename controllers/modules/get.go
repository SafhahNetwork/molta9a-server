package modules

import (
	"log"
	"molta9a/server/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetResponse struct {
	Module db.Module `json:"module"`
	NextID int
	PrevID int
}

func Get(c *gin.Context) {
	// Extract module ID from request parameters
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	// Query the database to fetch the module by ID
	var module db.Module
	result := db.DB.First(&module, moduleID)
	if result.Error != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var prevID int = -1
	if module.Index > 1 {
		result = db.DB.Model(&db.Module{}).Where("index = ?", module.Index-1).
			Order("index ASC").
			Limit(1).
			Pluck("id", &prevID)
		if err == gorm.ErrRecordNotFound {
			prevID = -1
		}
	}

	var nextID int = -1
	result = db.DB.Model(&db.Module{}).Where("index = ?", module.Index+1).
		Order("index ASC").
		Limit(1).
		Pluck("id", &nextID)
	if err == gorm.ErrRecordNotFound {
		nextID = -1
	}

	// Prepare the response
	response := GetResponse{
		Module: module,
		NextID: nextID,
		PrevID: prevID,
	}

	// Send the module as a response
	c.JSON(http.StatusOK, response)
}
