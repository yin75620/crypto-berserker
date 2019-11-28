package ftx

import (
	"fmt"
	"testing"
)

func TestBookerCenter(t *testing.T) {

	orderBookCenter := NewOrderBookCenter()
	market := "BTC/USD"
	if !orderBookCenter.IsExist(market) {
		orderBookCenter.Register(market)
	}

	booker := orderBookCenter.GetBooker(market)

	for {
		<-booker.UpdateChannel
		ask, bid := booker.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
