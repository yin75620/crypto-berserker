package ftx

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://ftx.com/ws/"

type FtxWebSocket struct {
}

func NewSocket() *FtxWebSocket {
	mws := &FtxWebSocket{}
	return mws
}

func (fws *FtxWebSocket) SubScribeOrderBook(marketName string) chan ob.OrderBookerResponseDetail {
	resChannel := make(chan ob.OrderBookerResponseDetail)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

func (fws *FtxWebSocket) doSubScribeOrderBook(marketName string, resChannel chan ob.OrderBookerResponseDetail) chan ob.OrderBookerResponseDetail {
	conn := createConn()
	market := marketName
	sendSubcribe(conn, "orderbook", market)

	go func() {

		for {
			response := ob.OrderBookerResponseDetail{}
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Println("error:", err)
				response.Error = err
				resChannel <- response
				fws.doSubScribeOrderBook(marketName, resChannel)
				return
			}
			//recv: {"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}
			//log.Printf("recv: %s", message)
			res := types.OrderBookSocketResponse{}

			json.Unmarshal(message, &res)
			response = toDetail(res)
			resChannel <- response
		}
	}()
	return resChannel
}

func toDetail(res types.OrderBookSocketResponse) ob.OrderBookerResponseDetail {
	obr := ob.OrderBookerResponseDetail{}
	if res.IsUpdate() {
		obr.Action = ob.Update
	} else {
		obr.Action = ob.Partial
	}

	obr.Asks = res.Data.Asks
	obr.Bids = res.Data.Bids
	obr.Checksum = res.Data.Checksum
	obr.Error = res.Error
	obr.Market = res.Market
	obr.Time = res.Data.Time

	return obr
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
