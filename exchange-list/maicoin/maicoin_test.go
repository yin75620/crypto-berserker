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
		Market:    "ethtwd",
		Side:      exc.Buy,
		Price:     5737.306,
		Size:      0.253332,
		OrderType: exc.MARKET,
		CoinPair:  exc.CoinPair{BaseCoin: "ETH", QuotedCoin: "TWD"},
	}
	mm.PostOrder(myOrder)

	//mm.GetAccountInfo()

	//mm.GetFill("usdttwd")

}
