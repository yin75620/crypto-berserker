package bybilinear

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

const WEBSOCKET_URL = "wss://stream.bybit.com/v5/public/linear"

type BybilinearWebSocket struct {
}

func NewSocket() *BybilinearWebSocket {
	mws := &BybilinearWebSocket{}
	return mws
}

func (fws *BybilinearWebSocket) SubScribeOrderBook(marketName string) chan ob.OrderBookerResponseDetail {
	resChannel := make(chan ob.OrderBookerResponseDetail)
	return fws.doSubScribeOrderBook(marketName, resChannel)
}

type OrderBookSocketResponseData struct {
	Symbol   string     `json:"s,omitempty"`
	Bids     [][]string `json:"b"`             // Bids 是出价，其中每个元素包含[价格, 数量]
	Asks     [][]string `json:"a,omitempty"`   // Asks 是要价，其中每个元素包含[价格, 数量]
	UpdateID int64      `json:"u,omitempty"`   // UpdateID 是更新ID，表示数据序列
	Sequence int64      `json:"seq,omitempty"` // Sequence 是交叉序列号
	Cts      int64      `json:"cts,omitempty"` // Cts 是匹配引擎产生此订单簿数据的时间戳
}

type OrderBookSocketResponse struct {
	Topic      string  `json:"topic,omitempty"`
	ActionType string  `json:"type,omitempty"`
	TimeStamp  float64 `json:"ts,omitempty"`
}

type SnapshotResponse struct {
	OrderBookSocketResponse
	Data OrderBookSocketResponseData `json:"data,omitempty"`
}

func (ssr *SnapshotResponse) toOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	res.Time = ssr.TimeStamp
	//res.checksum = ssr.Id
	res.Action = ob.Partial
	res.Market = ssr.Data.Symbol
	res.Asks = tool.TransToFloatTwoArray(ssr.Data.Asks)
	res.Bids = tool.TransToFloatTwoArray(ssr.Data.Bids)

	return res
}

type DeltaResponse struct {
	OrderBookSocketResponse
	Data OrderBookSocketResponseData `json:"data,omitempty"`
	Cts  int64                       `json:"cts"` // 匹配引擎时间戳
}

func (dr *DeltaResponse) toOrderBookDetail() ob.OrderBookerResponseDetail {
	res := ob.OrderBookerResponseDetail{}
	res.Time = dr.TimeStamp
	//res.checksum = dr.Id
	res.Action = ob.Update
	res.Market = dr.Data.Symbol

	res.Asks = tool.TransToFloatTwoArray(dr.Data.Asks)
	res.Bids = tool.TransToFloatTwoArray(dr.Data.Bids)

	return res
}

// /
// 為了錯誤重連，所以要再傳入 channel
func (fws *BybilinearWebSocket) doSubScribeOrderBook(marketName string, resChannel chan ob.OrderBookerResponseDetail) chan ob.OrderBookerResponseDetail {
	conn := createConn()
	market := marketName
	receive_res := sendSubcribe(conn, "orderbook", market)

	go func() {
		rep := addResultToChannel(receive_res, market)
		resChannel <- rep
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
			//recv: {"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}
			//log.Printf("recv: %s", message)

			firstResponse := OrderBookSocketResponse{}

			json.Unmarshal(message, &firstResponse)
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
			}
			resChannel <- response
		}
	}()
	return resChannel
}

func addResultToChannel(message []byte, market string) ob.OrderBookerResponseDetail {
	response := ob.OrderBookerResponseDetail{}
	firstResponse := OrderBookSocketResponse{}

	json.Unmarshal(message, &firstResponse)
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
	}
	return response
}

func createConn() *websocket.Conn {
	return tool.CreateConn(WEBSOCKET_URL)
}

type SubscriptionRequest struct {
	Args []string `json:"args,omitempty"` //subscribe,unsubscribe
	Op   string   `json:"op,omitempty"`   //subscribe,unsubscribe
}

func sendSubcribe(c *websocket.Conn, channelName string, marketName string) []byte {
	//"orderBookL2_25.BTCUSD"
	//ws.send('{"op": "subscribe", "args": ["orderBookL2_25.BTCUSD"]}');
	args := fmt.Sprintf("%s.50.%s", channelName, marketName)
	fmt.Println(args)
	req := SubscriptionRequest{
		Op:   "subscribe",
		Args: []string{args},
	}

	return tool.Send(c, req)

}
