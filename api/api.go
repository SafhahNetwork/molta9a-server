package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"molta9a/server/api/middleware"
	runnerAPI "molta9a/server/api/runner"

	authController "molta9a/server/controllers/auth"
	modulesController "molta9a/server/controllers/modules"
	tracksController "molta9a/server/controllers/tracks"
	userController "molta9a/server/controllers/user"
	"molta9a/server/db"
	"molta9a/server/redis"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Errors())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	/***
	    Auth
	***/

	auth := r.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	auth.POST("/reset-password", func(c *gin.Context) {
		c.String(http.StatusOK, "not implemented yet")
	})

	/***
	    Tracks
	***/
	tracks := r.Group("/tracks")
	tracks.GET("/get/:id", tracksController.Get)
	tracks.GET("/list", tracksController.List)

	tracks.GET("/list_progress", middleware.AuthMiddleware(), tracksController.ListProgress)
	tracks.GET("/get_progress/:id", middleware.AuthMiddleware(), tracksController.GetProgress)
	tracks.POST("/start/:id", middleware.AuthMiddleware(), tracksController.Start)

	/***
	    Modules
	***/
	modules := r.Group("/modules")
	modules.GET("/get/:id", modulesController.Get)
	modules.GET("/discussion", func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "not implemented yet")
	})
	modules.GET("/solutions", func(c *gin.Context) {
		c.String(http.StatusNotImplemented, "not implemented yet")
	})

	// modules.GET("/preview/:id", modulesController.Get)
	// modules.GET("/get/:id", middleware.AuthMiddleware(), modulesController.Get)

	/***
	   Exercise POST
	***/
	exercise := r.Group("/exercise")
	exercise.Use(middleware.AuthMiddleware())

	exercise.POST("/run/:id", runnerAPI.Run)
	exercise.POST("/submit/:id", runnerAPI.Submit)

	exercise.POST("/message", func(c *gin.Context) {
		c.String(http.StatusOK, "not implemented yet")
	})

	/***
	   User
	***/
	user := r.Group("/user")
	user.GET("/self", middleware.AuthMiddleware(), userController.PrivateGet)

	return r
}

func main() {
	db.Init()
	redis.Init()

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
