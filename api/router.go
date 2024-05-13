package api

import (
	"github.com/gin-gonic/gin"

	"molta9a/server/api/_middleware"
	authApi "molta9a/server/api/auth"
	coursesApi "molta9a/server/api/courses"
	exercisesApi "molta9a/server/api/exercises"
	selfApi "molta9a/server/api/self"
	userCoursesApi "molta9a/server/api/user-courses"
	userExercisesApi "molta9a/server/api/user-exercises"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Errors())

	/***
	    Auth
	***/

	auth := r.Group("/auth")
	auth.POST("/register", authApi.Register)
	auth.POST("/login", authApi.Login)

	/***
	    Courses
	***/
	courses := r.Group("/courses")

	courses.GET("/get/:id", coursesApi.Get)
	courses.GET("/list", coursesApi.List)

	/***
	    User Courses
	***/
	user_courses := r.Group("/user_courses")
	user_courses.Use(middleware.AuthMiddleware())

	user_courses.GET("/get/:course_id", userCoursesApi.Get)
	user_courses.POST("/start/:course_id", userCoursesApi.Start)
	user_courses.GET("/list", userCoursesApi.List)

	/***
	    Exercises
	***/
	exercises := r.Group("/exercises")
	exercises.GET("/get/:id", exercisesApi.Get)

	/***
	    User Exercises
	***/
	user_exercises := r.Group("/user_exercises")
	user_exercises.Use(middleware.AuthMiddleware())

	user_exercises.GET("/get/:exercise_id", userExercisesApi.Get)
	user_exercises.GET("/list/:course_id", userExercisesApi.List)

	user_exercises.POST("/run/:exercise_id", userExercisesApi.Run)
	user_exercises.POST("/submit/:exercise_id", userExercisesApi.Run) // TODO

	/***
	   Self
	***/
	user := r.Group("/self")
	user.GET("/get", middleware.AuthMiddleware(), selfApi.Get)

	return r
}
