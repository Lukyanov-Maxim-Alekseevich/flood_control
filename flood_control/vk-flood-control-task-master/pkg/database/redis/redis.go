package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     int
	Database int
	User     string
	Password string
}

type ManagedDatabase struct {
	RedisDb *redis.Client
}

func New(config *Config) (*ManagedDatabase, error) {
	redisUrl := fmt.Sprintf("redis://%s:%s@%s:%d/%d?protocol=3", config.User, config.Password, config.Host, config.Port, config.Database)

	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	redisDb := redis.NewClient(opts)

	managedDatabase := &ManagedDatabase{
		RedisDb: redisDb,
	}

	return managedDatabase, nil
}
