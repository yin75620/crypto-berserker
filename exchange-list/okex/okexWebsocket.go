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
				log.Println("error:", err)
				response.Error = err
				resChannel <- response
				fws.doSubScribeOrderBook(market, resChannel)
				return
			}

			unMessage, err := tool.FlateDecompress(message)
			if err != nil {
				log.Println("error:", err)
				response.Error = err
				resChannel <- response
				fws.doSubScribeOrderBook(market, resChannel)
				return
			}

			orderBookData := OrderBookData{}
			json.Unmarshal(unMessage, &orderBookData)
			//log.Println("orderBookData:", orderBookData)
			response = orderBookData.ToOrderBookDetail()
			response.Market = market

			resChannel <- response
		}
	}()
	return resChannel
}

func createConn(channelName string, marketName string) *websocket.Conn {
	//const levels = "5" //5,10,20
	//endpoint := fmt.Sprintf("%s/%s@%s%s", WEBSOCKET_URL, marketName, channelName, levels)
	return tool.CreateConn(WEBSOCKET_URL)
}

func sendSubcribe(c *websocket.Conn) {
	//"orderBookL2_25.BTCUSD"
	//ws.send('{"op": "subscribe", "args": ["orderBookL2_25.BTCUSD"]}');
	//endpoint := fmt.Sprintf("%s/%s@depth%s", baseURL, strings.ToLower(symbol), levels)
	//args := fmt.Sprintf("%sL2_25.%s", channelName, marketName)
	req := SubscriptionRequest{
		Op:   "subscribe",
		Args: []string{"swap/depth5:BTC-USD-SWAP"},
	}

	tool.SendFlate(c, req)
	//tool.SendPing(c, nil)
}
