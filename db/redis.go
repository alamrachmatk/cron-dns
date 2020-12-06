package db

import (
	"dns/config"

	"github.com/labstack/gommon/log"

	"github.com/garyburd/redigo/redis"
)

var Redis redis.Conn

func RedisInit() redis.Conn {
	conn, err := redis.Dial("tcp", config.RedisHost+":"+config.RedisPort)
	if err != nil {
		log.Fatal("Error connect to Redis server!")
	}
	return conn
}
