package binance

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://stream.binance.com:9443/ws"

type BinanceWebSocket struct {
}

func NewSocket() *BinanceWebSocket {
	mws := &BinanceWebSocket{}
	return mws
}

func (fws *BinanceWebSocket) SubScribeOrderBook(marketName string) chan ob.OrderBookerResponseDetail {
	resChannel := make(chan ob.OrderBookerResponseDetail)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

///
// 為了錯誤重連，所以要再傳入 channel
func (fws *BinanceWebSocket) doSubScribeOrderBook(marketName string, resChannel chan ob.OrderBookerResponseDetail) chan ob.OrderBookerResponseDetail {
	conn := createConn("depth", marketName)
	market := marketName
	sendSubcribe(conn)

	go func() {

		for {
			response := ob.OrderBookerResponseDetail{}

			_, message, err := conn.ReadMessage()

			//log.Println("message:", string(message))
			if err != nil {
				log.Println("error:", err)
				resChannel <- response
				fws.doSubScribeOrderBook(market, resChannel)
				return
			}
			//recv: {"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}
			//log.Printf("recv: %s", message)

			orderBookData := OrderBookData{}
			json.Unmarshal(message, &orderBookData)
			//log.Println("orderBookData:", orderBookData)
			response = orderBookData.ToOrderBookDetail()
			response.Market = market

			resChannel <- response
		}
	}()
	return resChannel
}

func createConn(channelName string, marketName string) *websocket.Conn {
	const levels = "5" //5,10,20
	endpoint := fmt.Sprintf("%s/%s@%s%s", WEBSOCKET_URL, marketName, channelName, levels)
	return tool.CreateConn(endpoint)
}

func sendSubcribe(c *websocket.Conn) {
	//"orderBookL2_25.BTCUSD"
	//ws.send('{"op": "subscribe", "args": ["orderBookL2_25.BTCUSD"]}');
	//endpoint := fmt.Sprintf("%s/%s@depth%s", baseURL, strings.ToLower(symbol), levels)
	//args := fmt.Sprintf("%sL2_25.%s", channelName, marketName)
	/*req := SubscriptionRequest{
		Op:   "subscribe",
		Args: []string{args},
	}*/

	//tool.Send(c, nil)
	tool.SendPing(c, nil)
}
