package order_booker

import (
	"fmt"
	"sync"
)

type OrderBookCenter struct {
	orderBookers map[string]*OrderBooker
	webSocket    IWebSocket
	mutex        sync.Mutex
}

func NewOrderBookCenter(websocket IWebSocket) *OrderBookCenter {
	obc := &OrderBookCenter{}
	obc.orderBookers = make(map[string]*OrderBooker)
	obc.webSocket = websocket
	return obc
}

func (obc *OrderBookCenter) IsExist(market string) bool {
	if _, ok := obc.orderBookers[market]; ok {
		return true
	}
	return false
}

func (obc *OrderBookCenter) GetBooker(market string) *OrderBooker {
	return obc.orderBookers[market]
}

// 給一個幣種，直接開始同步該幣種的最高最低價
func (obc *OrderBookCenter) Register(market string) (chan int, bool) {
	if _, ok := obc.orderBookers[market]; !ok {
		resChannel := obc.webSocket.SubScribeOrderBook(market)
		ob := NewOrderBooker(market, resChannel)

		obc.mutex.Lock()
		obc.orderBookers[market] = ob
		obc.mutex.Unlock()

		ob.Start()
		return ob.UpdateChannel, true
	} else {
		fmt.Println("Already Register:", market)
		return nil, false
	}
}
