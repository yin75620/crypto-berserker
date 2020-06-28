package okex

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://real.okex.com:8443/ws/v3"

type OkexWebSocket struct {
}

func NewSocket() *OkexWebSocket {
	mws := &OkexWebSocket{}
	return mws
}

func (fws *OkexWebSocket) SubScribeOrderBook(marketName string) chan ob.OrderBookerResponseDetail {
	resChannel := make(chan ob.OrderBookerResponseDetail)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

///
// 為了錯誤重連，所以要再傳入 channel
func (fws *OkexWebSocket) doSubScribeOrderBook(marketName string, resChannel chan ob.OrderBookerResponseDetail) chan ob.OrderBookerResponseDetail {
	conn := createConn("subscribe", marketName)

	market := marketName
	sendSubcribe(conn)

	go func() {

		for {
			response := ob.OrderBookerResponseDetail{}

			_, message, err := conn.ReadMessage()

			//log.Println("message:", string(message))
			if err != nil {
				log.Println("ReadMessage error:", err)
				response.Error = err
				resChannel <- response
				fws.doSubScribeOrderBook(market, resChannel)
				return
			}

			unMessage, err := tool.FlateDecompress(message)
			if err != nil {
				log.Println("Flate error:", err)
				response.Error = err
				resChannel <- response
				fws.doSubScribeOrderBook(market, resChannel)
				return
			}

			orderBookData := OrderBookData{}
			json.Unmarshal(unMessage, &orderBookData)
			response = orderBookData.ToOrderBookDetail()
			response.Market = market

			resChannel <- response
		}
	}()
	return resChannel
}

func createConn(channelName string, marketName string) *websocket.Conn {
	return tool.CreateConn(WEBSOCKET_URL)
}

func sendSubcribe(c *websocket.Conn) {
	req := SubscriptionRequest{
		Op:   "subscribe",
		Args: []string{"swap/depth5:BTC-USD-SWAP"},
	}

	tool.SendFlate(c, req)
}
