package ftx

import (
	"fmt"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

func TestSingleWebsocketInfo(t *testing.T) {

	socket := NewSingleSocket()
	coin := exc.CoinPair{"BTC", "USD"}
	ch := socket.SubScribeOrderBook(coin.GetMarketName())

	nob := ob.NewOrderBooker(coin.GetMarketName(), ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
func TestSingleWebsocketInfoFutures(t *testing.T) {

	socket := NewSingleSocket()
	marketName := "BTC-PERP"
	ch := socket.SubScribeOrderBook(marketName)

	nob := ob.NewOrderBooker(marketName, ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}

func TestSingleWebsocketMultiInfoFutures(t *testing.T) {

	socket := NewSingleSocket()
	marketName := "BTC-PERP"
	ch := socket.SubScribeOrderBook(marketName)
	nob := ob.NewOrderBooker(marketName, ch)

	marketName1 := "ETH-PERP"
	ch1 := socket.SubScribeOrderBook(marketName1)
	nob1 := ob.NewOrderBooker(marketName1, ch1)

	nob.Start()
	nob1.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)

		<-nob1.UpdateChannel
		ask1, bid1 := nob1.GetFirstPricePair()
		fmt.Println("ask1", ask1)
		fmt.Println("bid1", bid1)
	}
}
