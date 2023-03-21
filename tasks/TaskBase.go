package tasks

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"task/models"
	"task/redisClient"
	"time"
)

func HandleQueue(taskQueue string, method func(item interface{})) {
	var redisCon = redisClient.GetInstanceRedis()
	defer redisCon.Close()

	var popNum = 0
	var taskNum = 2 //默认2个进程  根据 waitNum 动态调整
	pool, _ := ants.NewPoolWithFunc(taskNum, func(i interface{}) {
		defer func() {
			if err := recover(); err != nil {
				models.CloseDb()
			}
		}()
		defer models.CloseDb()
		method(i)
	})
	for {
		res, err := redisCon.RPop(taskQueue).Result()
		if err != nil {
			//fmt.Println("pop err:", err)
		}
		if res == "" {
			fmt.Println("pop result nil")
			time.Sleep(3 * time.Second)
			continue
		}
		//fmt.Println("pop result :", res)

		_ = pool.Invoke(res)

		if popNum > 10 { //每pop 10个任务 查看一次任务队列排队数量  按排队数量动态控制协程数量
			waitNum, _ := redisCon.LLen(taskQueue).Result()
			newTaskNum := 2
			if waitNum < 10 {
				newTaskNum = 2
			} else if waitNum >= 10 && waitNum < 50 {
				newTaskNum = 5
			} else {
				newTaskNum = 10
			}
			if newTaskNum != taskNum {
				fmt.Printf("change task num from %d to %d\n", taskNum, newTaskNum)
				pool.Tune(newTaskNum)
				taskNum = newTaskNum
			}
			popNum = 0
		}

		popNum++
	}
}
