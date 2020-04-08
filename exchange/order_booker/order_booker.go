package order_booker

import (
	"sync"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type OrderBooker struct {
	MarketName            string
	lastTime              float64
	Asks                  map[float64]float64 //price,volume
	Bids                  map[float64]float64 //price,volume
	mutex                 sync.Mutex
	socketResponseChannel chan OrderBookerResponseDetail
	UpdateChannel         chan int
}

func NewOrderBooker(market string, resChannel chan OrderBookerResponseDetail) *OrderBooker {
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

		if res.Time < ob.lastTime {
			continue
		}
		ob.lastTime = res.Time

		ob.mutex.Lock()
		if res.IsUpdate() {
			updateToMap(res.Asks, &ob.Asks)
			updateToMap(res.Bids, &ob.Bids)
		} else {
			clearMap(&ob.Asks)
			clearMap(&ob.Bids)
			saveToMap(res.Asks, &ob.Asks)
			saveToMap(res.Bids, &ob.Bids)
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
