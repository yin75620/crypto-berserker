package maicoin

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/maicoin/types"
)

const WEBSOCKET_URL = "wss://max-ws.maicoin.com"

type MaincoinWebSocket struct {
	conn *websocket.Conn
}

func NewSocket() *MaincoinWebSocket {
	mws := &MaincoinWebSocket{}
	return mws
}

func (mws *MaincoinWebSocket) SubScribeOrderBook(market string) chan exc.OrderBookSocketResponse {
	conn := createConn()
	sendSubcribe(conn, "orderbook", market)

	resChannel := make(chan exc.OrderBookSocketResponse)
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			//recv: {"info":"error","msg":"unknown market btcusdts"}
			//recv: {"info":"subscribed","channel":"orderbook","market":"btcusdt"}
			//recv: {"info":"orderbook","timestamp":"1573130643672","action":"add","market":"btcusdt","id":75056203,"side":"buy","volume":"0.78","price":"9142.01","ord_type":"limit"}
			log.Printf("recv: %s", message)
			response := types.OrderBookSocketResponse{}
			json.Unmarshal(message, &response)
			resChannel <- response.OrderBookSocketResponse
		}
	}()
	return resChannel
}

func createConn() *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(WEBSOCKET_URL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func sendSubcribe(c *websocket.Conn, channelName string, marketName string) {
	req := types.SubscriptionRequest{
		Cmd:     "subscribe",
		Channel: channelName,
		Params: map[string]interface{}{
			"market": marketName,
		},
	}

	actionJson, err := json.Marshal(req)
	if err != nil {
		log.Fatal("json.Marshal:", err)
	}
	send(c, actionJson)
}

func send(c *websocket.Conn, json []byte) {
	err := c.WriteMessage(websocket.TextMessage, json)
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
