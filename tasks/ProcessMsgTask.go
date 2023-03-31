package tasks

import (
	"encoding/json"
	"fmt"
	"task/components"
)

func pushSendMsgQueue(toUserId string, message components.Message) {
	redis := components.GetInstanceRedis()
	data := struct {
		ToUserId string
		Message  components.Message
	}{
		ToUserId: toUserId,
		Message:  message,
	}
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("send msg err ", err)
	}
	redis.RPush("message:send:queue", string(b))
}

func ProcessMsg() {
	HandleQueue("message:process:queue", func(item interface{}) {
		fmt.Println("process msg ", item)
		message := components.Message{}
		err := json.Unmarshal([]byte(item.(string)), &message)
		if err != nil {
			fmt.Println("process msg err ", err)
		}
		switch message.Action {
		case components.ACTION_LOGIN:
			pushSendMsgQueue(message.UserId, message)
			return
		default:
			return
		}
	})
}
