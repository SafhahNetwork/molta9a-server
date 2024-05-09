#!/bin/bash
#
docker start postgres pgadmin registry zealous_franklin

wait 5

export POSTGRES_DSN="host=192.168.1.69 user=postgres password=admin dbname=mol port=5432 sslmode=disable"
export REDIS_URI="redis://@192.168.1.69:6379/0"

go run ./api
