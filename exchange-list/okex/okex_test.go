package okex

import (
	"fmt"
	"net/http"
	"testing"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

const (
	USD  = "USD"
	USDT = "USDT"
	BTC  = "BTC"
	FTT  = "FTT"
)

var m_OkexClient = NewOkex(http.DefaultClient)

func TestProducts(t *testing.T) {
	//res := m_OkexClient.GetProducts()
	//fmt.Println(fmt.Sprintf("res:%s", string(res)))
}

func TestAccount(t *testing.T) {
	res := m_OkexClient.GetAccountInfo()

	fmt.Println(string(res))
}
func TestGetPair(t *testing.T) {
	res1, res2 := m_OkexClient.GetAskBidPair(exc.CoinPair{BTC, USDT}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))
}

/*
func TestOrder(t *testing.T) {
	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		Market:    "FTT/USDT",
		Side:      exc.Buy,
		Price:     1.3,
		Size:      1,
		OrderType: exc.LIMIT,
	}
	m_OkexClient.PostOrder(myOrder)
}
*/
