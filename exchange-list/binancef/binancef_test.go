package binancef

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

var be = NewBinancef(http.DefaultClient)

func TestAccount(t *testing.T) {
	res := be.GetAccountInfo()

	fmt.Println(string(res))

	fmt.Println(be.GetAccount())
}

func TestOrder(t *testing.T) {

	start := time.Now()
	be.prepareMarketInfo()

	var myOrder exc.FuturesOrder = exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      exc.Buy,
			Price:     20.32,
			Size:      46.2,
			OrderType: exc.LIMIT,
		},
		Futures: exc.Futures{
			//ExpirationDate time.Time
			// 商品名
			TargetName: "SOL",
			// 計價貨幣類型
			QuoteCoin: "USDT",
		},
	}
	res, err := be.PostFuturesOrder(myOrder)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	fmt.Println("elapsed:", time.Since(start))
}

func TestCancelAllOrder(t *testing.T) {

	be.PostCancelAllOrder(exc.Futures{QuoteCoin: "USDT", TargetName: "SOL"})
}

func TestGetExchangeInfo(t *testing.T) {
	be.getExchangeInfo()
}

func TestPrepareLeverage(t *testing.T) {
	be.prepareLeverage()
}

func TestGetBalance(t *testing.T) {
	be.GetBalance()
}

func TestGetMaxOrderUSD(t *testing.T) {
	be.Prepare()
	fmt.Println(be.GetMaxOrderUSD("AXSUSDT"))
}

func TestPostLeverage(t *testing.T) {
	be.PostLeverage("BTCUSDT", 125)
}

func TestSetAllLeverage(t *testing.T) {
	be.setAllLeverage()
}

func TestUserTraders(t *testing.T) {

	userTrades := be.getUserTrades("GMTUSDT", time.Now().Add(-100*time.Hour), time.Now())

	for _, p := range userTrades {

		tool.PrintStructNameValue(p)
	}

}

func TestTightUserTraders(t *testing.T) {

	userTrades := be.GetTightUserTrades("GMTUSDT", time.Now().Add(-100*time.Hour), time.Now())
	fmt.Println(userTrades)
}

func TestFree(t *testing.T) {
	//be.Prepare()
	//be.getKline()
	//be.prepareLeverage()
	//be.getUserTrades()

}

/*
func TestSign(t *testing.T) {
	key := "vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A"
	secKey := "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"
	ce := NewBinance1(key, secKey)
	body := exc.JArray{
		"symbol":      "LTCBTC",
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    1,
		"price":       0.1,
		"recvWindow":  5000,
		"timestamp":   1499827319559,
	}
	ce.doRequest("GET", "Account", body)

	query := exc.JArray{
		"symbol":      "LTCBTC",
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
	}

	str := fmt.Sprintf("%s%s", query.ToValues().Encode(), body.ToValues().Encode())
	fmt.Println(str)

	//res, err := exc.GetParamHmacSHA256HexSign(secKey, "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559")
	res, err := exc.GetParamHmacSHA256HexSign(secKey, str)
	if err != nil {
		fmt.Println(err)
	}
	//symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559

	//res1 := "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71"
	fmt.Println(res)
}*/

func TestPricePair(t *testing.T) {
	ce := NewBinancef(http.DefaultClient)
	res1, res2 := ce.GetAskBidPair(exc.CoinPair{"BTC", "USDT"}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))

	//fmt.Println(string(res))
}

func TestGetWallet(t *testing.T) {

	ce := NewBinancef(http.DefaultClient)

	w := ce.GetWallet()
	fmt.Println(w)
}
