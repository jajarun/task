package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var userCons = make(map[string]map[*websocket.Conn]bool)

func handleMessage(msg string, userId string) {
	for conn, _ := range userCons[userId] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("收到%s的消息%s", userId, msg)))
		if err != nil {
			log.Println("write msg err:", err)
		}
	}
}

func Run() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		//value := r.URL.Query()
		//fmt.Println(value.Get("id"))
		s, _ := io.ReadAll(r.Body) //把	body 内容读入字符串 s
		a := string(s)
		fmt.Println("body ", a)
		_, _ = io.WriteString(w, "11111")
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(err)
			return
		}
		//defer conn.Close()
		query := request.URL.Query()
		userId := query.Get("user_id")
		if userId == "" {
			log.Println("user_id 为空 关闭连接")
			_ = conn.Close()
			return
		}
		if userCons[userId] == nil {
			userCons[userId] = make(map[*websocket.Conn]bool)
		}
		userCons[userId][conn] = true
		_ = conn.WriteMessage(websocket.TextMessage, []byte(userId+"连接成功"))
		log.Println(userId + "连接数" + strconv.Itoa(len(userCons[userId])))
		conn.SetCloseHandler(func(code int, text string) error {
			delete(userCons[userId], conn)
			log.Println(userId + "关闭连接")
			log.Println(userId + "连接数" + strconv.Itoa(len(userCons[userId])))
			return nil
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
			go handleMessage(string(msg), userId)

			//err = conn.WriteMessage(websocket.TextMessage, []byte("请再说一遍"))
			//if err != nil {
			//	log.Println("write msg err:", err)
			//}
		}
	})

	_ = http.ListenAndServe(":8181", nil)

}
