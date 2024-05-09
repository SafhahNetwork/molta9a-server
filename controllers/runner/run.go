package runner

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"molta9a/server/db"
	"molta9a/server/redis"
	"molta9a/server/runner/run"
	"molta9a/server/utils"
	"net/http"
	"time"
)

var ctx = context.Background()

func Run(userIP string, filesString string, moduleID uint, files []run.File) (output string, ok bool, err error, status int) {
	cacheKey := utils.GenKey([]string{"runner/run", filesString})

	// Check if the result is cached
	cachedResult, err := redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cached result found, return it
		log.Println("used cache YAAY with key: ", cacheKey)
		return cachedResult, true, nil, http.StatusOK
	}
	log.Println("no cache :< for key: ", cacheKey, " - ", err)

	// Query the database to fetch the module by ID
	var module db.Module
	result := db.DB.First(&module, moduleID)
	if result.Error != nil || *module.Exercise == "" {
		log.Println(err)
		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
	}

	exerciseJson, err := base64.StdEncoding.DecodeString(*module.Exercise)
	if err != nil {
		log.Println(err)
		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
	}

	var exercise db.Exercise
	err = json.Unmarshal(exerciseJson, &exercise)
	if err != nil {
		log.Println(err)
		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
	}

	image := exercise.DockerImage
	env := exercise.Env
	args := exercise.Args

	// Run container
	b, ok, err := run.Run(userIP, image, env, args, files)
	message := string(b)
	if err != nil {
		message = err.Error()
	}

	// Cache the result
	redis.Client.Set(ctx, cacheKey, message, 2*time.Hour)

	return message, true, nil, http.StatusOK
}
