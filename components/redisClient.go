package components

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"sync"
)

var instanceRedis *redis.Client
var onceRedis sync.Once

func GetInstanceRedis() *redis.Client {
	onceRedis.Do(func() {
		addr := viper.Get("redis.addr")
		redisCon := redis.NewClient(&redis.Options{
			Addr: addr.(string),
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
