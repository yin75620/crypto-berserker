package CrossExchange

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yin75620/crypto-berserker/message_tool"

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

	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"tester"})
	bybit := bybit.NewBybit(http.DefaultClient)

	crossPairsTable := map[string][]CrossPair{}
	cp := CrossPair{
		askExchange:  ft,
		bidExchange:  bybit,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  1.0,
	}

	cpArray := []CrossPair{
		cp,
	}
	crossPairsTable[cp.GetName()] = cpArray

	matchMap := map[string]CrossPair{}

	matchCp := CrossPair{
		askExchange:  bybit,
		bidExchange:  ft,
		askPricePair: exc.PricePair{9500, 1000},
		bidPricePair: exc.PricePair{9700, 1000},
		orderVolume:  1.0,
	}
	matchMap[matchCp.GetName()] = matchCp

	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USD",
	}
	message_tool.StartTelegram()
	init := *NewCrossExchangeInit()
	isClose, res := positionCloseCheck(crossPairsTable, matchMap, futures, init)

	assert.True(t, isClose)
	assert.Zero(t, len(res))
}
