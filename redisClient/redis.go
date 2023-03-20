package redisClient

import (
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

var instanceRedis *redis.Client
var once sync.Once

func GetInstanceRedis() *redis.Client {
	once.Do(func() {
		redisCon := redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
		})
		result, err := redisCon.Ping().Result()
		if err != nil {
			fmt.Println("ping err:", err)
		}
		fmt.Println(result)
		instanceRedis = redisCon
	})
	return instanceRedis
}
