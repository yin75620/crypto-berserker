package tool

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func CreateConn(url string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func Send(c *websocket.Conn, req interface{}) {
	json, err := json.Marshal(req)
	if err != nil {
		log.Fatal("json.Marshal:", err)
		return
	}

	err = c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println("send error:", err)
		return
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("receive: %s\n", msg)
}
