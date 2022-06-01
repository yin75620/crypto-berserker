package ftx

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

var mftx = NewFtx(http.DefaultClient,
	FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})

func TestHistoryInfo(t *testing.T) {
	now := time.Now()
	sec := now.Unix()

	prices, _ := mftx.GetHistoricalPrices("BTC/USD", 60, 10, sec-300, sec)
	log.Println(prices)
}

func TestInfo(t *testing.T) {

	fmt.Println(string(mftx.GetAccountInfo()))

	fmt.Println(mftx.GetAccount())

	//ftx.GetMarkets()
}

func TestOrder(t *testing.T) {

	// var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
	// 	Market:    "FTT/USD",
	// 	Side:      exc.Buy,
	// 	Price:     1.1,
	// 	Size:      1.1255,
	// 	OrderType: exc.LIMIT,
	// }
	//mftx.PostOrder(myOrder)

}

func TestFuturesOrder(t *testing.T) {
	order := exc.FuturesOrder{}
	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USDT",
	}

	order.Futures = futures
	order.Size = 0.0001
	order.Price = 9201.1205
	order.Side = exc.Buy
	order.OrderType = exc.LIMIT
	mftx.PostFuturesOrder(order)
}

func TestMarket(t *testing.T) {
	fmt.Println(string(mftx.GetMarkets()))
}

func TestBalance(t *testing.T) {

}

func TestAPI(t *testing.T) {
	res := mftx.GetLendingInfo()
	li := res.GetLendInfo("USD")
	fmt.Print(res)

	lo := NewLendOrder("USD", li.Lendable, li.MinRate)
	mftx.PostLendingOffer(*lo)
}

func TestDeleteAllOrder(t *testing.T) {
	order := FtxOrder{}
	order.Market = "FIDA/USD"
	order.Side = "buy"
	order.Price = 0.5
	order.Size = 1
	order.OrderType = LIMIT
	mftx.doPostOrder(order)
	res, err := mftx.DeleteAllOrders()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
