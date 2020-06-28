package okex

import (
	"fmt"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

func TestWebsocketInfo(t *testing.T) {

	socket := NewSocket()
	futures := exc.Futures{time.Time{}, "BTC", "USDT"}
	ch := socket.SubScribeOrderBook(futures.GetSwapNameUpper())

	nob := ob.NewOrderBooker(futures.GetSwapNameUpper(), ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}
