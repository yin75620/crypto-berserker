package CrossExchange

import (
	"net/http"
	"testing"

	"github.com/yin75620/crypto-berserker/exchange-list/bybit"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func createTestCrossExchange() *CrossExchange {
	exchanges := []exc.Exchange{}
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	bybit := bybit.NewBybit(http.DefaultClient)
	exchanges = append(exchanges, ft)
	exchanges = append(exchanges, bybit)

	ce := NewCrossExchange(exchanges)
	return ce
}

func TestMain(t *testing.T) {
	ce := createTestCrossExchange()
	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USD",
	}
	ce.SetFuturesArray([]exc.Futures{futures})
	ce.Start()
}

func TestOrder(t *testing.T) {
	/*
		ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			"tester"})
		futures := exc.Futures{
			//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
			TargetName: "BTC",
			QuoteCoin:  "USD",
		}
		executeOrder(ft, futures, 8000, exc.Ask, 1)*/
}

func TestPositionCloseCheck(t *testing.T) {

	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USD",
	}

	crossPairMap := map[string]CrossPair{}
	crossPairsTable := map[string][]CrossPair{}
	positionCloseCheck(crossPairsTable, crossPairMap, futures)
}
