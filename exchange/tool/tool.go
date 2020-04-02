package tool

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

const RetryCount = 10

func CreateConn(url string) *websocket.Conn {
	var c *websocket.Conn
	var err error

	for i := 0; i < RetryCount; i++ {
		c, err = doCreateConn(url)
		if err != nil {
			continue
		}
		return c
	}

	log.Fatal(err)
	return c
}

func doCreateConn(url string) (*websocket.Conn, error) {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("dial:", err)
		return c, err
	}
	return c, nil
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

func SendPing(c *websocket.Conn, req interface{}) {
	err := c.WriteMessage(websocket.PingMessage, []byte{})
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

func TransToFloatTwoArray(askStrArrays []interface{}) [][]float64 {
	res := [][]float64{}
	for _, array := range askStrArrays {
		askFloatArray := []float64{}
		sArray := array.([]interface{})
		for _, s := range sArray {
			res, _ := strconv.ParseFloat(s.(string), 64)
			askFloatArray = append(askFloatArray, res)
		}
		res = append(res, askFloatArray)
	}
	return res
}
