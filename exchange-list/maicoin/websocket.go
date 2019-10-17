package maicoin

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/maicoin/types"
)

const WEBSOCKET_URL = "wss://max-ws.maicoin.com"

func Start() {
	c, _, err := websocket.DefaultDialer.Dial(WEBSOCKET_URL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//defer c.Close()

	market := "btcusdt"

	req := types.SubscriptionRequest{
		Cmd:     "subscribe",
		Channel: "orderbook",
		Params: map[string]interface{}{
			"market": market,
		},
	}

	actionJson, _ := json.Marshal(req)
	log.Println(string(actionJson))
	send(c, actionJson)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			//recevier(message)
		}
	}()
}

func send(c *websocket.Conn, json []byte) {
	err := c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println(err)
		return
	}
	/*_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("receive: %s\n", msg)*/
}
