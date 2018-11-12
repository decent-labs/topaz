package redis

import (
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

// Cli is used to communicate with Redis throughout our application
type Cli struct {
	conn redis.Conn
}

var cliInstance *Cli

// Connect created a new connection to redis and returns an instance of type Cli
func Connect() (conn *Cli) {
	if cliInstance == nil {
		cliInstance = new(Cli)
		var err error

		cliInstance.conn, err = redis.Dial("tcp", fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")))
		if err != nil {
			panic(err)
		}

		if _, err := cliInstance.conn.Do("AUTH", os.Getenv("REDIS_PASSWORD")); err != nil {
			cliInstance.conn.Close()
			panic(err)
		}
	}

	return cliInstance
}

// SetValue takes data and configuration to store a value in redis
func (cli *Cli) SetValue(key string, value string, expiration ...interface{}) error {
	_, err := cli.conn.Do("SET", key, value)

	if err == nil && expiration != nil {
		cli.conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

// GetValue returns the value stored at a specific key
func (cli *Cli) GetValue(key string) (interface{}, error) {
	return cli.conn.Do("GET", key)
}
