package bybilinear

import (
	"encoding/json"
	"fmt"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	ob "github.com/yin75620/crypto-berserker/exchange/order_booker"
)

func TestWebsocketInfo(t *testing.T) {

	socket := NewSocket()
	coin := exc.CoinPair{"BTC", "USDT"}
	ch := socket.SubScribeOrderBook(coin.GetLinkMakertNameUpper())

	nob := ob.NewOrderBooker(coin.GetLinkMakertNameUpper(), ch)
	nob.Start()
	for {
		<-nob.UpdateChannel
		ask, bid := nob.GetFirstPricePair()
		fmt.Println("ask", ask)
		fmt.Println("bid", bid)
	}
}

func TestConvert(t *testing.T) {
	message := []byte("{\"topic\":\"orderbook.50.BTCUSDT\",\"type\":\"delta\",\"ts\":1709972689094,\"data\":{\"s\":\"BTCUSDT\",\"b\":[],\"a\":[[\"68410.10\",\"0.308\"],[\"68418.20\",\"0.010\"]],\"u\":16044212,\"seq\":139157701116},\"cts\":1709972689090}")
	fmt.Println(string(message))
	deltaResponse := DeltaResponse{}
	json.Unmarshal(message, &deltaResponse)
	fmt.Println(deltaResponse)

}
