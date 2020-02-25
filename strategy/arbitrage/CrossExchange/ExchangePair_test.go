package CrossExchange

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func TestGetMaxHoldVolume(t *testing.T) {
	start := time.Now()
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	ft.GetAccountInfo()
	bybit := bybit.NewBybit(http.DefaultClient)
	bybit.GetAccountInfo()

	cp := CrossPair{
		askExchange:  ft,
		bidExchange:  bybit,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  200,
	}

	fmt.Println(cp.GetMaxHoldVolume())
	fmt.Println(time.Since(start).Seconds())
}
