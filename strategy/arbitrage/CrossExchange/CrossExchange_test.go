package CrossExchange

import (
	"net/http"
	"testing"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func TestMain(t *testing.T) {
	exchanges := []exc.Exchange{}
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	exchanges = append(exchanges, ft)
	ce := NewCrossExchange(exchanges)
	futures := exc.Futures{
		ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName:     "BTC",
		QuoteCoin:      "USD",
	}
	ce.setFuturesArray([]exc.Futures{futures})
	ce.Start()
}
