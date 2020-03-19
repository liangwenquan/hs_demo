package redisLib

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"fmt"
)

var (
	client *redis.Client
)

func NewClient()  {
	client = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})

	pong, err := client.Ping().Result()
	fmt.Println("初始化redis:", pong, err)
}

func GetClient() (c *redis.Client) {

	return client
}