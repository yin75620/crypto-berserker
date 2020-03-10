package binance

import (
	"fmt"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

func TestWebsocketInfo(t *testing.T) {

	socket := NewSocket()
	coin := exc.CoinPair{"BTC", "USDT"}
	ch := socket.SubScribeOrderBook(coin.GetLinkMakertName())

	nob := ob.NewOrderBooker(coin.GetLinkMakertName(), ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
