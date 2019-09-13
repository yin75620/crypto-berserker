package ftxotc

import (
	"fmt"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/setting"
)

var mftexchange = NewFtxotc(http.DefaultClient,
	FtxotcInit{
		setting.FTXOTC_KEY,
		setting.FTXOTC_SECRET_KEY,
	})

func TestInfo(t *testing.T) {

	//fmt.Println(string(mftexchange.GetBalance()))

	fmt.Println(string(mftexchange.GetQuote(exc.CoinPair{"XRP", "USDT"}, "sell", 10)))

	//fmt.Println(string(mftexchange.GetAllTradingPair()))
	//fmt.Println(string(mftexchange.GetQuoteByID(52864844)))
}

func TestOrder(t *testing.T) {

	// var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
	// 	Market:    "FTT/USD",
	// 	Side:      exc.Buy,
	// 	Price:     1.1,
	// 	Size:      1.1255,
	// 	OrderType: exc.LIMIT,
	// }
	//mftexchange.PostOrder(myOrder)

}

func TestMarket(t *testing.T) {
	//fmt.Println(string(mftexchange.GetMarkets()))
}
