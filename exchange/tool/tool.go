package tool

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const RetryCount = 8640
const WaitSecond = 10

func CreateConn(url string) *websocket.Conn {
	var c *websocket.Conn
	var err error

	for i := 0; i < RetryCount; i++ {
		c, err = doCreateConn(url)
		if err != nil {
			waitChannel := make(chan int)
			go func() {
				time.Sleep(time.Second * WaitSecond)
				waitChannel <- 0
			}()
			<-waitChannel
			log.Println("Error: retryCreateConn, count:", i)
			continue
		}
		return c
	}

	log.Println("Error: CreateConn Fail")
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

func Send(c *websocket.Conn, req interface{}) []byte {

	msg, err := doSend(c, req)
	if err != nil {
		log.Fatal("send has error:", err)
		return []byte{}
	}
	log.Printf("receive: %s\n", msg)
	return msg
}

func SendFlate(c *websocket.Conn, req interface{}) {
	msg, err := doSend(c, req)
	if err != nil {
		log.Fatal("SendFlate has error:", err)
		return
	}
	unMsg, err := FlateDecompress(msg)
	if err != nil {
		log.Fatal("FlateDecompress has error:", err)
		return
	}
	log.Printf("SendFlate receive: %s\n", unMsg)
}

func doSend(c *websocket.Conn, req interface{}) ([]byte, error) {
	json, err := json.Marshal(req)
	if err != nil {
		log.Fatal("json.Marshal:", err)
		return []byte{}, err
	}

	err = c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println("send error:", err)
		return []byte{}, err
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return []byte{}, err
	}
	return msg, nil
}

func FlateDecompress(data []byte) ([]byte, error) {
	return ioutil.ReadAll(flate.NewReader(bytes.NewReader(data)))
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

func TransInterfaceToFloatTwoArray(askStrArrays []interface{}) [][]float64 {
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

func TransToFloatTwoArray(askStrArrays [][]string) [][]float64 {
	res := [][]float64{}
	for _, array := range askStrArrays {
		askFloatArray := []float64{}
		sArray := array
		for _, s := range sArray {
			res, _ := strconv.ParseFloat(s, 64)
			askFloatArray = append(askFloatArray, res)
		}
		res = append(res, askFloatArray)
	}
	return res
}
