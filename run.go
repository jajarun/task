package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"task/redisClient"
	"task/tasks"
	"time"
)

func printRedisPool(stats *redis.PoolStats) {
	fmt.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n",
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

func printRedisOption(opt *redis.Options) {
	fmt.Printf("Network=%v\n", opt.Network)
	fmt.Printf("Addr=%v\n", opt.Addr)
	fmt.Printf("Password=%v\n", opt.Password)
	fmt.Printf("DB=%v\n", opt.DB)
	fmt.Printf("MaxRetries=%v\n", opt.MaxRetries)
	fmt.Printf("MinRetryBackoff=%v\n", opt.MinRetryBackoff)
	fmt.Printf("MaxRetryBackoff=%v\n", opt.MaxRetryBackoff)
	fmt.Printf("DialTimeout=%v\n", opt.DialTimeout)
	fmt.Printf("ReadTimeout=%v\n", opt.ReadTimeout)
	fmt.Printf("WriteTimeout=%v\n", opt.WriteTimeout)
	fmt.Printf("PoolSize=%v\n", opt.PoolSize)
	fmt.Printf("MinIdleConns=%v\n", opt.MinIdleConns)
	fmt.Printf("MaxConnAge=%v\n", opt.MaxConnAge)
	fmt.Printf("PoolTimeout=%v\n", opt.PoolTimeout)
	fmt.Printf("IdleTimeout=%v\n", opt.IdleTimeout)
	fmt.Printf("IdleCheckFrequency=%v\n", opt.IdleCheckFrequency)
	fmt.Printf("TLSConfig=%v\n", opt.TLSConfig)

}

func main() {
	redisCon := redisClient.GetInstanceRedis()
	for i := 0; i < 100; i++ {
		redisCon.LPush("test:task", i)
	}

	tasks.HandleQueueTimed("test:task", func(item interface{}) {
		fmt.Println("method process:", item)
		time.Sleep(time.Second)
	})
}
