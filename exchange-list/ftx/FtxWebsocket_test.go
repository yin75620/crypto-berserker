package ftx

import (
	"fmt"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func TestWebsocketInfo(t *testing.T) {

	socket := NewSocket()
	coin := exc.CoinPair{"BTC", "USD"}
	ch := socket.SubScribeOrderBook(coin.GetMarketName())

	nob := NewOrderBooker(coin.GetMarketName(), ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
func TestWebsocketInfoFutures(t *testing.T) {

	socket := NewSocket()
	marketName := "BTC-PERP"
	ch := socket.SubScribeOrderBook(marketName)

	nob := NewOrderBooker(marketName, ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
