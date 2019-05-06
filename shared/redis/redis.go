package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

var (
	pool        *redis.Pool
	redisServer string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("couldn't load dotenv:", err.Error())
	}

	redisServer = fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisServer)
		},
	}
}

// SetValue takes data and configuration to store a value in redis
func SetValue(key string, value interface{}, expiration ...interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)

	if err == nil && expiration != nil {
		conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

// GetValue returns the value stored at a specific key
func GetValue(key string) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	return conn.Do("GET", key)
}
