package bybit

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://stream-testnet.bybit.com/realtime"

type ByBitWebSocket struct {
	//conn *websocket.Conn
	Asks map[float64]float64 //price,volume
	Bids map[float64]float64 //price,volume
}

func NewSocket() *ByBitWebSocket {
	mws := &ByBitWebSocket{}
	return mws
}

func (fws *ByBitWebSocket) SubScribeOrderBook(marketName string) chan types.OrderBookSocketResponse {
	resChannel := make(chan types.OrderBookSocketResponse)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

從這裡繼續，準備對 Response 做轉換
//https://github.com/bybit-exchange/bybit-official-api-docs/blob/master/en/websocket.md#orderBook25_v2
/*
//snapshot type format. The data is ordered by price, from buy to sell
{
     "topic":"orderBookL2_25.BTCUSD",
     "type":"snapshot",
     "data":[
        {
            "price":"2999.00",
            "symbol":"BTCUSD",
            "id":29990000,
            "side":"Buy",
            "size":9
        },
        {
            "price":"3001.00",
            "symbol":"BTCUSD",
            "id":30010000,
            "side":"Sell",
            "size":10
        }
     ],
     "cross_seq":11518,
     "timestamp_e6":1555647164875373
}
*/

type OrderBookSocketResponseData struct {
	Time     float64     `json:"time,omitempty"`
	Checksum int64       `json:"checksum,omitempty"`
	Asks     [][]float64 `json:"asks,omitempty"`
	Bids     [][]float64 `json:"bids,omitempty"`
	Action   string      `json:"action,omitempty"`
}

type OrderBookSocketResponse struct {
	Topic      string                      `json:"topic,omitempty"`
	ActionType ActionType                  `json:"action,omitempty"`
	Data       OrderBookSocketResponseData `json:"data,omitempty"`
	//Success
}

func (fws *ByBitWebSocket) doSubScribeOrderBook(marketName string, resChannel chan types.OrderBookSocketResponse) chan types.OrderBookSocketResponse {
	conn := createConn()
	market := marketName
	sendSubcribe(conn, "orderbook", market)

	go func() {

		for {
			response := types.OrderBookSocketResponse{}
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Println("error:", err)
				resChannel <- response
				fws.doSubScribeOrderBook(marketName, resChannel)
				return
			}
			//recv: {"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}
			//log.Printf("recv: %s", message)

			json.Unmarshal(message, &response)
			resChannel <- response
		}
	}()
	return resChannel
}

func createConn() *websocket.Conn {
	return tool.CreateConn(WEBSOCKET_URL)
}

type SubscriptionRequest struct {
	Args []string `json:"args,omitempty"` //subscribe,unsubscribe
	Op   string   `json:"op,omitempty"`   //subscribe,unsubscribe
}

func sendSubcribe(c *websocket.Conn, channelName string, marketName string) {
	//"orderBookL2_25.BTCUSD"
	args := fmt.Sprintf("%sL2_25.%s", channelName, marketName)
	req := SubscriptionRequest{
		Op:   "subscribe",
		Args: {args},
	}

	tool.Send(c, req)
}
