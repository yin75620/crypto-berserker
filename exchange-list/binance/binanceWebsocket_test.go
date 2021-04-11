package binance

import (
	"fmt"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

func TestWebsocketInfo(t *testing.T) {

	socket := NewSocket()
	// 如果是不存在的交易對，會停住
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
