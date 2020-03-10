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
	//conn *websocket.Conn
	Asks map[float64]float64 //price,volume
	Bids map[float64]float64 //price,volume
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
type OrderBookSocketResponseData struct {
	Price  float64 `json:"price,string,omitempty"`
	Symbol string  `json:"symbol,omitempty"`
	Id     int64   `json:"id,omitempty"`
	Side   string  `json:"side,omitempty"`
	Size   float64 `json:"size,omitempty"`
}

type OrderBookSocketResponse struct {
	Topic      string  `json:"topic,omitempty"`
	ActionType string  `json:"type,omitempty"`
	CrossSeq   int64   `json:"cross_seq,omitempty"`
	TimeStamp  float64 `json:"timestamp_e6,omitempty"`
}

type SnapshotResponse struct {
	OrderBookSocketResponse
	Data []OrderBookSocketResponseData `json:"data,omitempty"`
}

func (ssr *SnapshotResponse) toOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	res.Time = ssr.TimeStamp
	//res.checksum = ssr.Id
	res.Action = ob.Partial
	//res.Market = ssr.Symbol
	transToAskBid(&res.Asks, &res.Bids, ssr.Data)
	return res
}

type DeltaDetail struct {
	Delete []OrderBookSocketResponseData `json:"delete,omitempty"`
	Update []OrderBookSocketResponseData `json:"update,omitempty"`
	Insert []OrderBookSocketResponseData `json:"insert,omitempty"`
}

type DeltaResponse struct {
	OrderBookSocketResponse
	Data DeltaDetail `json:"data,omitempty"`
}

func (dr *DeltaResponse) toOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	res.Time = dr.TimeStamp
	//res.checksum = dr.Id
	res.Action = ob.Update
	//res.Market = dr.Symbol

	transToAskBid(&res.Asks, &res.Bids, dr.Data.Delete)
	transToAskBid(&res.Asks, &res.Bids, dr.Data.Update)
	transToAskBid(&res.Asks, &res.Bids, dr.Data.Insert)
	return res
}

func transToAskBid(asks, bids *[][]float64, data []OrderBookSocketResponseData) {
	for _, pItem := range data {
		if pItem.Side == "Buy" {
			*bids = append(*bids, []float64{pItem.Price, pItem.Size})
		} else if pItem.Side == "Sell" {
			*asks = append(*asks, []float64{pItem.Price, pItem.Size})
		} else {
			log.Println("unknow pItem.Side:", pItem.Side)
		}
	}
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

			log.Println("message:", string(message))
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
			log.Println(orderBookData)
			response = orderBookData.ToOrderBookDetail()
			response.Market = market

			/*json.Unmarshal(message, &firstResponse)
			if firstResponse.ActionType == "snapshot" {
				snapshotResponse := SnapshotResponse{}
				json.Unmarshal(message, &snapshotResponse)
				response = snapshotResponse.toOrderBookDetail()
				response.Market = market
			} else if firstResponse.ActionType == "delta" {
				deltaResponse := DeltaResponse{}
				json.Unmarshal(message, &deltaResponse)
				response = deltaResponse.toOrderBookDetail()
				response.Market = market
			}*/
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
