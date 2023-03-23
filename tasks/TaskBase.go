package tasks

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
	"task/redisClient"
	"time"
)

func HandleQueueTimed(taskQueue string, method func(item interface{})) {
	var channelQueue = make(chan string, 100)
	redisCon := redisClient.GetInstanceRedis()
	queueNum, _ := redisCon.LLen(taskQueue).Result()
	newTaskNum := 2
	if queueNum < 10 {
		newTaskNum = 2
	} else if queueNum >= 10 && queueNum < 50 {
		newTaskNum = 5
	} else {
		newTaskNum = 10
	}
	fmt.Printf("master task start task num:%d \n", newTaskNum)
	newTaskNum += 1 //还需要开启一个线程用于写入通道
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(newTaskNum)
	go func(chan string) {
		defer waitGroup.Done()
		for {
			item, _ := redisCon.RPop(taskQueue).Result()
			if item == "" {
				fmt.Println("pop result nil")
				break
			}
			fmt.Println("pop result:", item)
			channelQueue <- item
		}
		close(channelQueue)
	}(channelQueue)

	for i := 1; i < newTaskNum; i++ {
		go func(chan string) {
			defer waitGroup.Done()
			for {
				item, ok := <-channelQueue
				if !ok {
					break
				}
				fmt.Println("sub task get result:", item)
				method(item)
			}
		}(channelQueue)
	}
	waitGroup.Wait()
	fmt.Println("master task end")
}

func HandleQueue(taskQueue string, method func(item interface{})) {
	var redisCon = redisClient.GetInstanceRedis()
	defer redisCon.Close()

	var popNum = 0
	var taskNum = 2 //默认2个进程  根据 waitNum 动态调整
	pool, _ := ants.NewPoolWithFunc(taskNum, method)
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
