package runnerAPI

import (
	"encoding/json"
	"log"
	"molta9a/server/controllers/runner"
	"molta9a/server/runner/run"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	moduleIDStr := c.Param("id")
	moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	// Read body
	var files []run.File
	if err := c.BindJSON(&files); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid files"})
		return
	}
	filesString, err := json.Marshal(files)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid files"})
		return
	}

	// Get data from request context
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	output, ok, err, status := runner.Run(ip, string(filesString), uint(moduleID), files)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Respond with JSON
	c.JSON(status, gin.H{"Output": output, "Ok": ok})
}
