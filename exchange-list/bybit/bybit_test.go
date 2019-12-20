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
