package userExercisesController

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"log"
// 	"molta9a/server/db"
// 	"molta9a/server/redis"
// 	"molta9a/server/runner/run"
// 	"molta9a/server/utils"
// 	"net/http"
// 	"time"
//
// 	"gorm.io/gorm"
// )
//
// func Submit(user *db.User, userIP string, filesString string, moduleID uint, files []run.File) (output string, ok bool, err error, status int) {
// 	cacheKey := utils.GenKey([]string{"runner/submit", filesString})
//
// 	// Check if the result is cached
// 	cachedResult, err := redis.Client.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		// Cached result found, return it
// 		log.Println("used cache YAAY with key: ", cacheKey)
// 		return cachedResult, false, nil, http.StatusOK
// 	}
// 	log.Println("no cache :< for key: ", cacheKey, " - ", err)
//
// 	// Query the database to fetch the module by ID
// 	var userModule db.UserModule
// 	result := db.DB.First(&userModule).Where("module_id = ? AND user_id = ?", moduleID, user.ID)
// 	if userModule.Status == "pass" {
// 		return "", false, utils.ErrModuleFinished, http.StatusNotAcceptable
// 	}
//
// 	// Query the database to fetch the module by ID
// 	var module db.Module
// 	result = db.DB.First(&module, moduleID)
// 	if result.Error != nil || *module.Exercise == "" {
// 		log.Println(err)
// 		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
// 	}
//
// 	exerciseJson, err := base64.StdEncoding.DecodeString(*module.Exercise)
// 	if err != nil {
// 		log.Println(err)
// 		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
// 	}
//
// 	var exercise db.Exercise
// 	err = json.Unmarshal(exerciseJson, &exercise)
// 	if err != nil {
// 		log.Println(err)
// 		return "", false, utils.ErrModuleNotFound, http.StatusNotFound
// 	}
//
// 	image := exercise.DockerImage
// 	env := exercise.Env
// 	args := exercise.Args
// 	env = append(env, "TEST_EXERCISE=true")
//
// 	// Run container
// 	b, ok, err := run.Run(userIP, image, env, args, files)
// 	message := string(b)
// 	if err != nil {
// 		message = err.Error()
// 	}
//
// 	// If user pass the test update the DB
// 	if ok {
// 		// create UserModule
// 		userModule = db.UserModule{
// 			UserID:   user.ID,
// 			TrackID:  module.TrackID,
// 			QuestID:  module.QuestID,
// 			ModuleID: module.ID,
// 			TotalXp:  module.TotalXp,
// 			Status:   "pass",
// 			Files:    &filesString,
// 			Output:   &message,
// 		}
// 		result = db.DB.Save(&userModule)
// 		if result.Error != nil {
// 			log.Println(err)
// 			return message, ok, utils.ErrSomethingWentWrong, http.StatusInternalServerError
// 		}
//
// 		// add XP and other stats to user
// 		updates := map[string]interface{}{
// 			"total_xp":        gorm.Expr("total_xp + ?", module.TotalXp),
// 			"total_exercises": gorm.Expr("total_exercises + ?", 1),
// 		}
// 		result = db.DB.Model(db.User{}).Where("id = ?", user.ID).Updates(updates)
// 		if result.Error != nil {
// 			log.Println(err)
// 			return message, ok, utils.ErrSomethingWentWrong, http.StatusInternalServerError
// 		}
// 	} else {
// 		// Cache the result
// 		redis.Client.Set(ctx, cacheKey, message, 2*time.Hour)
// 	}
//
// 	return message, ok, nil, http.StatusOK
// }
