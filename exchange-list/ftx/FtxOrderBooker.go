package ftx

import (
	"fmt"
	"sync"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx/types"
)

type OrderBooker struct {
	MarketName            string
	lastTime              float64
	Asks                  map[float64]float64 //price,volume
	Bids                  map[float64]float64 //price,volume
	mutex                 sync.Mutex
	socketResponseChannel chan types.OrderBookSocketResponse
	UpdateChannel         chan int
}

func NewOrderBooker(market string, resChannel chan types.OrderBookSocketResponse) *OrderBooker {
	ob := &OrderBooker{}
	ob.Asks = map[float64]float64{}
	ob.Bids = map[float64]float64{}
	ob.UpdateChannel = make(chan int)
	ob.MarketName = market
	ob.socketResponseChannel = resChannel
	return ob
}

func (ob *OrderBooker) Start() {
	go ob.Receive()
}

func (ob *OrderBooker) Receive() {
	for {
		res := <-ob.socketResponseChannel

		if res.Error != nil {
			ob.mutex.Lock()
			clearMap(&ob.Asks)
			clearMap(&ob.Bids)
			ob.mutex.Unlock()
			continue
		}

		// 比較市場
		if res.Market != ob.MarketName {
			continue
		}

		if res.Data.Time < ob.lastTime {
			continue
		}
		ob.lastTime = res.Data.Time

		ob.mutex.Lock()
		if res.IsUpdate() {
			updateToMap(res.Data.Asks, &ob.Asks)
			updateToMap(res.Data.Bids, &ob.Bids)
		} else {
			clearMap(&ob.Asks)
			clearMap(&ob.Bids)
			saveToMap(res.Data.Asks, &ob.Asks)
			saveToMap(res.Data.Bids, &ob.Bids)
		}
		ob.mutex.Unlock()
		ob.UpdateChannel <- 0
	}
}

//GetFirstPricePair return Ask, Bid
func (ob *OrderBooker) GetFirstPricePair() (exc.PricePair, exc.PricePair) {

	ob.mutex.Lock()
	minKey, minValue := mapMin(ob.Asks)
	maxKey, maxValue := mapMax(ob.Bids)
	ob.mutex.Unlock()

	return exc.PricePair{Price: minKey, Volume: minValue},
		exc.PricePair{Price: maxKey, Volume: maxValue}
}

func mapMax(numbers map[float64]float64) (float64, float64) {
	if len(numbers) == 0 {
		return 0, 0
	}
	var maxNumber float64
	for maxNumber = range numbers {
		break
	}
	for n := range numbers {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber, numbers[maxNumber]
}

func mapMin(numbers map[float64]float64) (float64, float64) {
	if len(numbers) == 0 {
		return 0, 0
	}
	var minNumber float64
	for minNumber = range numbers {
		break
	}
	for n := range numbers {
		if n < minNumber {
			minNumber = n
		}
	}
	return minNumber, numbers[minNumber]
}

func updateToMap(priceArray [][]float64, myMap *map[float64]float64) {
	for _, v := range priceArray {
		index := v[0]
		value := v[1]
		if value == 0 {
			if _, ok := (*myMap)[index]; ok {
				delete((*myMap), index)
			}
		} else {
			(*myMap)[index] = value
		}

	}
}

func saveToMap(priceArray [][]float64, myMap *map[float64]float64) {
	for _, v := range priceArray {
		index := v[0]
		value := v[1]
		if _, ok := (*myMap)[index]; ok {
			(*myMap)[index] += value
		} else {
			(*myMap)[index] = value
		}
	}
}

func clearMap(myMap *map[float64]float64) {
	*myMap = map[float64]float64{}
}

type OrderBookCenter struct {
	orderBookers map[string]*OrderBooker
	webSocket    *FtxWebSocket
	mutex        sync.Mutex
}

func NewOrderBookCenter() *OrderBookCenter {
	obc := &OrderBookCenter{}
	obc.orderBookers = make(map[string]*OrderBooker)
	obc.webSocket = NewSocket()
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
