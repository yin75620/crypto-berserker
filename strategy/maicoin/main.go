package main

import (
	"fmt"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/maicoin"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mMai = maicoin.NewMaicoin(http.DefaultClient)
var mTri = Tri.NewTriangular(mMai)

func main() {
	fmt.Println("TEST")

	/*mTri.SetCoinArrays([][]string{
		[]string{"BTC", "TWD"},
		[]string{"BTC", "USDT", "TWD"},
	})*/
	//mTri.Start()

	var mm = maicoin.NewMaicoin(http.DefaultClient)
	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		Market:    "btcusdt",
		Side:      exc.Buy,
		Price:     31,
		Size:      1,
		OrderType: exc.LIMIT,
		CoinPair:  exc.CoinPair{BaseCoin: "USDT", QuotedCoin: "TWD"},
	}
	mm.PostOrder(myOrder)

	//mm.GetAccountInfo()

	//mm.GetFill("usdttwd")

}
