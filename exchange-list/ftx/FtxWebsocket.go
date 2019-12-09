package ftx

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://ftx.com/ws/"

type FtxWebSocket struct {
	//conn *websocket.Conn
	Asks map[float64]float64 //price,volume
	Bids map[float64]float64 //price,volume
}

func NewSocket() *FtxWebSocket {
	mws := &FtxWebSocket{}
	return mws
}

func (fws *FtxWebSocket) SubScribeOrderBook(marketName string) chan types.OrderBookSocketResponse {
	resChannel := make(chan types.OrderBookSocketResponse)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

func (fws *FtxWebSocket) doSubScribeOrderBook(marketName string, resChannel chan types.OrderBookSocketResponse) chan types.OrderBookSocketResponse {
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

func sendSubcribe(c *websocket.Conn, channelName string, marketName string) {
	req := types.SubscriptionRequest{
		Op:      "subscribe",
		Channel: channelName,
		Market:  marketName,
	}

	tool.Send(c, req)
}
