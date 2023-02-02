package bybilinear

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange/tool"
)

var defaultFutures = exc.Futures{
	TargetName: "AXS",
	QuoteCoin:  "USDT",
}

var bb = NewBybilinear(http.DefaultClient)

func TestAccount(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)
	res := ce.GetAccountInfo()

	fmt.Println(string(res))

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

	bb.prepareMarketInfo()

	var myOrder exc.FuturesOrder = exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      "buy",
			Price:     22.11111,
			Size:      1.00162621642642,
			OrderType: exc.LIMIT,
		},
		Futures: defaultFutures,
	}
	res, _ := bb.PostFuturesOrder(myOrder)
	fmt.Println(res)
}

func TestOrderCancelAll(t *testing.T) {
	ce := NewBybilinear(http.DefaultClient)

	ce.PostCancelAllOrder(defaultFutures)
}

func TestGetWallet(t *testing.T) {

	ce := NewBybilinear(http.DefaultClient)

	w := ce.GetWallet()
	fmt.Println(w)
}

func TestGetInstrumentInfo(t *testing.T) {
	bb.prepareMarketInfo()
	ins := bb.getInstrumentInfo()
	symbolName := "SOLUSDT"
	for _, value := range ins.Results.List {
		if value.Symbol == symbolName {
			println(symbolName)
			fmt.Println(value.LotSizeFilter.QtyStep)
			fmt.Println(value.PriceFilter.TickSize)

			assert.Equal(t, value.LotSizeFilter.QtyStep, bb.marketInfos[symbolName].VolumeIncrement, "qty not equal")
			assert.Equal(t, value.PriceFilter.TickSize, bb.marketInfos[symbolName].PriceIncrement, "price not equal")
		}
	}

}

func TestPrepareLeverage(t *testing.T) {
	bb.prepareLeverage()
	fmt.Println(bb.LeverageInfos["BTCUSDT"])
}

func TestSetAllLeverage(t *testing.T) {
	bb.setAllLeverage()
}

func TestPostLeverage(t *testing.T) {
	bb.PostSetLeverage("BTCUSDT", 100)
}

func TestUserTrade(t *testing.T) {
	userTrades := bb.getUserTrades("APTUSDT", time.Time{}, time.Time{})
	for _, v := range userTrades {
		tool.PrintStructNameValue(v)
	}

	assert.NotEqual(t, len(userTrades), 0)
}

func TestTightUserTrade(t *testing.T) {
	tightUserTrades := bb.GetTightUserTrades("STGUSDT")
	fmt.Println(tightUserTrades)

	assert.NotEqual(t, len(tightUserTrades), 0)
}

func TestFree(t *testing.T) {
	bb.Prepare()
	fmt.Println(bb.GetMaxOrderUSD("BTCUSDT"))
}
