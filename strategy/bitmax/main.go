package main

import (
	"fmt"
	"net/http"

	bitmax "github.com/yin75620/crypto-berserker/bitmax"
	exc "github.com/yin75620/crypto-berserker/exchange"
)

const (
	USD  = "USD"
	USDT = "USDT"
	BTC  = "BTC"
	FTT  = "FTT"
)

var m_bitmaxClient = bitmax.NewBitmax(http.DefaultClient)

func main() {
	res := m_bitmaxClient.GetAccountInfo()

	fmt.Println(string(res))

	res1, res2 := m_bitmaxClient.GetAskBidPair(exc.CoinPair{FTT, USDT}, 1)
	fmt.Println(res1)
	fmt.Println(res2)
}
