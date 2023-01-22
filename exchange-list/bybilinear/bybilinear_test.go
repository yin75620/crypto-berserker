package bybilinear

import (
	"fmt"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func TestAccount(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)
	res := ce.GetAccountInfo()

	fmt.Println(string(res))

	fmt.Println(ce.GetAccount())
}

func TestPricePair(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)
	res1, res2 := ce.GetAskBidPair(exc.CoinPair{"BTC", "USDT"}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))

	//fmt.Println(string(res))
}

func TestOrderExecute(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)

	var myOrder exc.FuturesOrder = exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      "buy",
			Price:     7700.0112640126,
			Size:      0.00162621642642,
			OrderType: exc.LIMIT,
		},
		Futures: exc.Futures{
			//ExpirationDate time.Time
			// 商品名
			TargetName: "BTC",
			// 計價貨幣類型
			QuoteCoin: "USDT",
		},
	}
	ce.PostFuturesOrder(myOrder)
}

func TestOrderCancelAll(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)

	ce.PostCancelAllOrder(exc.Futures{
		TargetName: "BTC",
		QuoteCoin:  "USDT",
	})
}

func TestGetWallet(t *testing.T) {

	ce := NewBybilinear(http.DefaultClient)

	w := ce.GetWallet()
	fmt.Println(w)
}
