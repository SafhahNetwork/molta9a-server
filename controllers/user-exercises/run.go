package userExercisesController

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"molta9a/server/db"
	"molta9a/server/redis"
	"molta9a/server/runner"
	"molta9a/server/utils"
	"net/http"
	"time"
)

var ctx = context.Background()

func Run(userIP string, filesString string, moduleID uint, files []utils.File) (output string, ok bool, err error, status int) {
	cacheKey := utils.GenKey([]string{"runner/run", filesString})

	// Check if the result is cached
	// cachedResult, err := redis.Client.Get(ctx, cacheKey).Result()
	// if err == nil {
	// 	// Cached result found, return it
	// 	log.Println("used cache YAAY with key: ", cacheKey)
	// 	return cachedResult, true, nil, http.StatusOK
	// }
	// log.Println("no cache :< for key: ", cacheKey, " - ", err)

	// Query the database to fetch the module by ID
	var exercise db.Exercise
	result := db.DB.First(&exercise, moduleID)
	if result.Error != nil || *exercise.ExerciseConfig == "" {
		log.Println(result.Error)
		return "", false, utils.ErrExerciseNotFound, http.StatusNotFound
	}

	exerciseJson, err := base64.StdEncoding.DecodeString(*exercise.ExerciseConfig)
	if err != nil {
		log.Println(err)
		return "", false, utils.ErrExerciseNotFound, http.StatusNotFound
	}

	var exerciseCfg db.ExerciseConfig
	err = json.Unmarshal(exerciseJson, &exercise)
	if err != nil {
		log.Println(err)
		return "", false, utils.ErrExerciseNotFound, http.StatusNotFound
	}

	image := exerciseCfg.DockerImage
	env := exerciseCfg.Env
	args := exerciseCfg.Args

	log.Println("ARGGGGGGGGGGG ", args)

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
