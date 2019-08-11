package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	ftx "github.com/yin75620/crypto-berserker/ftx"
	bsk "github.com/yin75620/crypto-berserker/setting"
)

type ActionRequest struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type LoginRequest struct {
	Op   string             `json:"op"`
	Args LoginRequestDetail `json:"args"`
}

type LoginRequestDetail struct {
	Key  string `json:"key"`
	Time int64  `json:"time"` // integer current timestamp (in milliseconds)
	Sign string `json:"sign"` //SHA256 HMAC of the following string, using your API secret: <time>websocket_login
	//Subaccount string `json:"subaccount"` // (optional) subaccount name
}

func main() {

	ts := ftx.GetTimeSpan()
	tsStr := ftx.GetTimeSpanStr(ts)
	tsSign, _ := ftx.GetParamHmacSHA256HexSign(bsk.FTX_API_SECRET_KEY, tsStr+ftx.WEBSOCKET_LOGIN_KEY_WORD)
	log.Println(tsStr + ftx.WEBSOCKET_LOGIN_KEY_WORD)

	loginRequest := LoginRequest{
		Op: "login",
		Args: LoginRequestDetail{
			Key:  bsk.FTX_KEY,
			Time: ts,
			Sign: tsSign,
		},
	}

	loginJson, _ := json.Marshal(loginRequest)
	log.Println(string(loginJson))

	c, _, err := websocket.DefaultDialer.Dial("wss://ftexchange.com/ws/", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	actionRequest := ActionRequest{
		Op:      "subscribe",
		Channel: "orderbook",
		Market:  "FTT/USD",
	}
	actionJson, _ := json.Marshal(actionRequest)
	//log.Println(actionJson)
	send(c, actionJson)

	//log.Println(loginJson)
	//send(c, loginJson)

	////開 gorountine 收訊息
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
		}
	}()
	////

	///開啟伺服器讓程式留著
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	///
}

func send(c *websocket.Conn, json []byte) {
	err := c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Println(err)
		return
	}
	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("receive: %s\n", msg)
}
