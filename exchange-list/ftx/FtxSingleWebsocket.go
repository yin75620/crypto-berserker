package ftx

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

type FtxSingleWebSocket struct {
	conn           *websocket.Conn
	subscribeNames map[string](chan ob.OrderBookerResponseDetail)
}

func NewSingleSocket() *FtxSingleWebSocket {
	fsws := &FtxSingleWebSocket{}
	fsws.subscribeNames = make(map[string](chan ob.OrderBookerResponseDetail))
	return fsws
}

// 同一個名字只允許註冊一次
func (fsws *FtxSingleWebSocket) SubScribeOrderBook(marketName string) chan ob.OrderBookerResponseDetail {

	resChannel := make(chan ob.OrderBookerResponseDetail)
	if value, ok := fsws.subscribeNames[marketName]; ok {
		resChannel = value
	} else {
		fsws.subscribeNames[marketName] = resChannel
	}

	fsws.doSubScribeOrderBook(marketName)
	return resChannel
}

func (fsws *FtxSingleWebSocket) doSubScribeAllOrderBook() {
	if fsws.conn == nil {
		fsws.conn = createConn()
	}

	for k := range fsws.subscribeNames {
		sendSubcribe(fsws.conn, "orderbook", k)
	}
}

func (fsws *FtxSingleWebSocket) doSubScribeOrderBook(marketName string) {
	isFirst := fsws.conn == nil
	if isFirst {
		fsws.conn = createConn()

	}
	market := marketName
	sendSubcribe(fsws.conn, "orderbook", market)

	if isFirst {
		fsws.startLoop()
	}
}

func (fsws *FtxSingleWebSocket) startLoop() {
	go func() {

		for {
			response := ob.OrderBookerResponseDetail{}
			_, message, err := fsws.conn.ReadMessage()

			if err != nil {
				log.Println("error:", err)
				response.Error = err
				fsws.conn = nil
				fsws.sendDetailToAll(response)
				fsws.doSubScribeAllOrderBook()
				return
			}
			//recv: {"channel": "orderbook", "market": "BTC/USD", "type": "update", "data": {"time": 1574775055.2251372, "checksum": 3250722053, "bids": [], "asks": [[7089.5, 0.0], [7124.5, 47.3641]], "action": "update"}}
			//log.Printf("recv: %s", message)
			res := types.OrderBookSocketResponse{}

			json.Unmarshal(message, &res)

			response = toDetail(res)
			fsws.subscribeNames[res.Market] <- response
		}
	}()
}

func (fsws *FtxSingleWebSocket) sendDetailToAll(response ob.OrderBookerResponseDetail) {
	for _, v := range fsws.subscribeNames {
		v <- response
	}
}
