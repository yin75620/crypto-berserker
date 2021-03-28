package bitmax

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

var m_bitmaxClient = NewBitmax(http.DefaultClient)

func TestProducts(t *testing.T) {
	res := m_bitmaxClient.GetProducts()
	fmt.Println(fmt.Sprintf("res:%s", string(res)))
}

func TestAccount(t *testing.T) {
	res := m_bitmaxClient.GetAccountInfo()

	fmt.Println(string(res))

	res1, res2 := m_bitmaxClient.GetAskBidPair(exc.CoinPair{FTT, USDT}, 1)
	fmt.Println(fmt.Sprintf("res1:%g", res1.Price))
	fmt.Println(fmt.Sprintf("res1:%g", res1.Volume))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Price))
	fmt.Println(fmt.Sprintf("res2:%g", res2.Volume))

	// var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
	// 	Market:    "FTT/USDT",
	// 	Side:      exc.Buy,
	// 	Price:     1.3,
	// 	Size:      1,
	// 	OrderType: exc.LIMIT,
	// }
	// m_bitmaxClient.PostOrder(myOrder)
}

func TestOrder(t *testing.T) {
	order := exc.ExchangeOrder{}
	order.CoinPair = exc.CoinPair{"FTT", "USDT"}
	order.Side = "buy"
	order.Price = 32
	order.Size = 1
	order.OrderType = exc.LIMIT

	res, err := m_bitmaxClient.PostOrder(order)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
