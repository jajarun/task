package route

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"task/components"
	"task/controller/acc"
	"task/tasks"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var accRouteHandles = make(map[string]func(msg interface{}))

func setAccHandle(action string, handleFunc func(msg interface{})) {
	accRouteHandles[action] = handleFunc
}

func initAccHandles() {
	setAccHandle(components.ACTION_LOGIN, acc.Login)
}

func AccRouteInit() {
	initAccHandles()
	components.InitClientManager()
	http.HandleFunc("/ws", handleAccCon)

	go handleMessage()
	go tasks.ProcessMsg()
}

func handleMessage() {
	tasks.HandleQueue("message:send:queue", func(item interface{}) {
		fmt.Println("message send ", item)
		data := struct {
			ToUserId string
			Message  components.Message
		}{}
		if json.Unmarshal([]byte(item.(string)), &data) == nil {
			components.SendMsg(data.ToUserId, data.Message)
		}
	})
}

func getServer() string {
	return "server_1"
}

func handleAccCon(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	query := request.URL.Query()
	userId := query.Get("user_id")
	if userId == "" {
		log.Println("user_id 为空 关闭连接")
		_ = conn.Close()
		return
	}
	conn.SetCloseHandler(func(code int, text string) error {
		components.PutMessage(&components.ClientMsg{
			Conn:   conn,
			UserId: userId,
			Msg: components.Message{
				UserId: userId,
				Action: components.ACTION_LOGOUT,
				Server: getServer(),
			},
		})
		return nil
	})
	components.PutMessage(&components.ClientMsg{
		Conn:   conn,
		UserId: userId,
		Msg: components.Message{
			UserId: userId,
			Action: components.ACTION_LOGIN,
			Server: getServer(),
		},
	})
	for {
		messageType, msg, err := conn.ReadMessage()
		log.Println("msg type:", strconv.Itoa(messageType))
		if err != nil {
			_ = conn.Close()
			log.Println("read msg err:", err)
			return
		}
		log.Println("read msg:" + string(msg))
		message := components.Message{}
		if json.Unmarshal(msg, &message) == nil {
			message.UserId = userId
			message.Server = getServer()
			components.PutMessage(&components.ClientMsg{
				Conn:   conn,
				UserId: userId,
				Msg:    message,
			})
		}
	}
}
