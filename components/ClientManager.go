package components

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

const (
	ACTION_LOGIN  = "login"
	ACTION_LOGOUT = "logout"
)

type ClientMsg struct {
	Conn   *websocket.Conn
	UserId string
	Msg    Message
}

type Message struct {
	UserId string
	Action string
	Data   interface{}
	Server string
}

var clientManager ClientManager

func InitClientManager() {
	clientManager = ClientManager{
		userConns:    make(map[string]*websocket.Conn),
		messageQueue: make(chan Message, 1000),
	}
}

func PutMessage(clientMsg *ClientMsg) {
	switch clientMsg.Msg.Action {
	case ACTION_LOGIN:
		clientManager.login(clientMsg)
	case ACTION_LOGOUT:
		clientManager.logout(clientMsg)
	default:
	}
	pushProcessMsg(clientMsg.Msg)
}

func pushProcessMsg(message Message) {
	redis := GetInstanceRedis()
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println("push message err ", err)
	}
	redis.RPush("message:process:queue", string(b))
}

type ClientManager struct {
	userConns    map[string]*websocket.Conn
	messageQueue chan Message
	lock         sync.Mutex
}

func (cm *ClientManager) login(clientMsg *ClientMsg) {
	cm.lock.Lock()
	_, ok := cm.userConns[clientMsg.UserId]
	if !ok {
		cm.userConns[clientMsg.UserId] = clientMsg.Conn
	}
	cm.lock.Unlock()
}

func (cm *ClientManager) logout(clientMsg *ClientMsg) {
	cm.lock.Lock()
	_, ok := cm.userConns[clientMsg.UserId]
	if ok {
		delete(cm.userConns, clientMsg.UserId)
	}
	cm.lock.Unlock()
}

func SendMsg(toUserId string, message Message) {
	conn, ok := clientManager.userConns[toUserId]
	if !ok {
		return
	}
	b, err := json.Marshal(message)
	if err != nil {
		fmt.Println(toUserId, " send msg err", err)
	}
	_ = conn.WriteMessage(websocket.TextMessage, b)
}
