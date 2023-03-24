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
		query := request.URL.RawQuery
		_ = conn.WriteMessage(websocket.TextMessage, []byte(query+"连接成功"))
		conn.SetCloseHandler(func(code int, text string) error {
			log.Println(query + "关闭连接")
			return nil
		})
		go func(conn *websocket.Conn) {
			for {
				messageType, msg, err := conn.ReadMessage()
				log.Println("msg type:", strconv.Itoa(messageType))
				if err != nil {
					_ = conn.Close()
					log.Println("read msg err:", err)
					return
				}
				log.Println("read msg:" + string(msg))

				err = conn.WriteMessage(websocket.TextMessage, []byte("请再说一遍"))
				if err != nil {
					log.Println("write msg err:", err)
				}
			}
		}(conn)
	})

	_ = http.ListenAndServe(":8181", nil)

}
