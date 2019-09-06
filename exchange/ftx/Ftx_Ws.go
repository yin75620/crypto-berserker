package ftx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"

	bsk "github.com/yin75620/crypto-berserker/setting"
)

type ActionRequest struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type FillsRequest struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
}

type LoginRequest struct {
	Args LoginRequestDetail `json:"args"`
	Op   string             `json:"op"`
}

type LoginRequestDetail struct {
	Key  string `json:"key"`
	Sign string `json:"sign"` //SHA256 HMAC of the following string, using your API secret: <time>websocket_login
	Time int64  `json:"time"` // integer current timestamp (in milliseconds)
	//Subaccount string `json:"subaccount"` // (optional) subaccount name
}

const (
	//WEBSOCKET_URL = "wss://ftexchange.com/ws/"
	WEBSOCKET_URL = "wss://ftx.com/ws/"
)

type Receiver func([]byte)

func getLoginRequest() LoginRequest {
	var ts int64 = exc.GetTimeSpan()
	tsStr := exc.GetTimeSpanStr(ts)
	finalStr := tsStr + "websocket_login"
	log.Println("time+websocket:" + finalStr)
	tsSign, _ := exc.GetParamHmacSHA256HexSign(setting.FTX_API_SECRET_KEY, finalStr)
	log.Println("sign:" + tsSign)

	loginRequest := LoginRequest{
		Op: "login",
		Args: LoginRequestDetail{
			Key:  bsk.FTX_KEY,
			Sign: tsSign,
			Time: ts,
		},
	}
	return loginRequest
}

func sendRequest(actionObj interface{}, recevier Receiver) {

	c, _, err := websocket.DefaultDialer.Dial(WEBSOCKET_URL, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	//defer c.Close()

	actionJson, _ := json.Marshal(actionObj)
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
			recevier(message)
		}
	}()
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

func WsTest() {
	/*
		actionRequest := ActionRequest{
			Op:      "subscribe",
			Channel: "orderbook",
			Market:  "BTC/USD",
		}
		sendRequest(actionRequest, func(recv []byte) {
			log.Printf("actionRecv: %s", recv)
		})
	*/
	loginReq := getLoginRequest()
	sendRequest(loginReq, func(recv []byte) {
		log.Printf(": %s", recv)
	})
	/*
		fillsRequest := FillsRequest{
			Op:      "subcribe",
			Channel: "fills",
		}
		sendRequest(fillsRequest, func(recv []byte) {
			log.Printf("fillsRequestRecv: %s", recv)
		})*/

	///開啟伺服器讓程式留著
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":18080", nil))
	///
}
