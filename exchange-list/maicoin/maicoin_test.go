package maicoin

import (
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

var mMai = NewMaicoin(http.DefaultClient)
var mTri = Tri.NewTriangular(mMai)

func TestOrder(t *testing.T) {

	/*mTri.SetCoinArrays([][]string{
		[]string{"BTC", "TWD"},
		[]string{"BTC", "USDT", "TWD"},
	})*/
	//mTri.Start()

	var mm = NewMaicoin(http.DefaultClient)
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
