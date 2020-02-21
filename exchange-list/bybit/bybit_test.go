package bybit

import (
	"fmt"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

func TestAccount(t *testing.T) {
	ce := NewBybit(http.DefaultClient)
	res := ce.GetAccountInfo()

	fmt.Println(string(res))
}

func TestPricePair(t *testing.T) {
	ce := NewBybit(http.DefaultClient)
	res1, res2 := ce.GetAskBidPair(exc.CoinPair{"BTC", "USDT"}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))

	//fmt.Println(string(res))
}

func TestOrderExecute(t *testing.T) {
	ce := NewBybit(http.DefaultClient)

	var myOrder exc.FuturesOrder = exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      "buy",
			Price:     8900,
			Size:      1,
			OrderType: exc.MARKET,
		},
		Futures: exc.Futures{
			//ExpirationDate time.Time
			// 商品名
			TargetName: "BTC",
			// 計價貨幣類型
			QuoteCoin: "USD",
		},
	}
	ce.PostFuturesOrder(myOrder)
	// var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
	// 	Market:    "FTT/USD",
	// 	Side:      exc.Buy,
	// 	Price:     1.1,
	// 	Size:      1.1255,
	// 	OrderType: exc.LIMIT,
	// }
	//mftx.PostOrder(myOrder)
}

func TestGetWallet(t *testing.T) {

	ce := NewBybit(http.DefaultClient)

	w := ce.GetWallet()
	fmt.Println(w)
}
