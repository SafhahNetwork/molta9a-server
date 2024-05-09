package redis

import (
	"github.com/redis/go-redis/v9"
	"os"
)

var Client *redis.Client

func Init() {
	var uri = os.Getenv("REDIS_URI")

	opt, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	Client = redis.NewClient(opt)
}
