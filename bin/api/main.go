package main

import (
	"molta9a/server/api"
	"molta9a/server/db"
	"molta9a/server/redis"
)

func main() {
	db.Init()
	redis.Init()

	r := api.SetupRouter()
	r.Run(":8080")
}
