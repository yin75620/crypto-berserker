package CrossExchange

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybilinear"
)

func TestGetMaxHoldVolume(t *testing.T) {
	start := time.Now()

	bnf := binancef.NewBinancef(http.DefaultClient)
	bnf.GetAccountInfo()
	bybilinear := bybilinear.NewBybilinear(http.DefaultClient)
	bybilinear.GetAccountInfo()

	cp := CrossPair{
		Symbol:       "BTCUSDT",
		askExchange:  bnf,
		bidExchange:  bybilinear,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  200,
	}

	fmt.Println(cp.GetMaxHoldVolume())
	fmt.Println(time.Since(start).Seconds())
}
