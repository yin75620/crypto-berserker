package CrossExchange

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybilinear"
	"github.com/yin75620/crypto-berserker/jmath"
	"github.com/yin75620/crypto-berserker/message_tool"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

var bnf = binancef.NewBinancef(http.DefaultClient)
var bbl = bybilinear.NewBybilinear(http.DefaultClient)
var exchanges = []exc.Exchange{bnf, bbl}
var ce = NewCrossExchange(exchanges)

func TestMain(t *testing.T) {
	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USD",
	}
	ce.SetFuturesArray([]exc.Futures{futures})
	ce.Start()
}

func TestExchangesTrade(t *testing.T) {
	exchanges := []exc.Exchange{}
	bybilinear := bybilinear.NewBybilinear(http.DefaultClient)
	bnf := binancef.NewBinancef(http.DefaultClient)

	exchanges = append(exchanges, bybilinear)
	exchanges = append(exchanges, bnf)

	ce := NewCrossExchange(exchanges)

	futures := exc.Futures{
		TargetName: "BTC",
		QuoteCoin:  "USDT",
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

func TestTelegramPost(t *testing.T) {

	//message_tool.StartTelegram()
	//message_tool.SendBroadcastArcherGroup("abc: 123\nchange: desk\nabcdefghijklmnopqrstuvwxyz: 2345678911")
	//fmt.Println(fmt.Sprintf("time: %s", time.Now().UTC()))

	exchanges := []exc.Exchange{}
	bbl := bybilinear.NewBybilinear(http.DefaultClient)
	bnf := binancef.NewBinancef(http.DefaultClient)

	exchanges = append(exchanges, bbl)
	exchanges = append(exchanges, bnf)

	ce := NewCrossExchange(exchanges)

	futures := exc.Futures{
		TargetName: "BTC",
		QuoteCoin:  "USDT",
	}

	first := bbl
	second := bnf

	ftxBit := CrossPair{
		askExchange:  first,
		bidExchange:  second,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  10.0,
	}

	bybitFtx := CrossPair{
		askExchange:  second,
		bidExchange:  first,
		askPricePair: exc.PricePair{9500, 1000},
		bidPricePair: exc.PricePair{9700, 1000},
		orderVolume:  10.0,
	}

	cpArray := []CrossPair{
		ftxBit,
		bybitFtx,
	}

	matchMap := getCrossPairMap(ce.exchanges, futures)
	matchCrossPair := getMatchCrossPair("positionPairName", cpArray, matchMap)

	message_tool.StartTelegram()
	message_tool.SendBroadcastArcherGroup(matchCrossPair.GetProfitString())

	/*
		start := time.Now()
		time.Sleep(time.Second)
		elapsed := time.Since(start)

		fmt.Println(fmt.Sprintf("time: %s\r\nTime elapsed: %v", time.Now().UTC(), elapsed))*/
}

func TestPositionCloseCheck(t *testing.T) {
	/*
		first := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			"tester"})
		bybit := bybit.NewBybit(http.DefaultClient)
	*/
	first := bybilinear.NewBybilinear(http.DefaultClient)
	second := binancef.NewBinancef(http.DefaultClient)

	crossPairsTable := map[string][]CrossPair{}
	ftxBit := CrossPair{
		askExchange:  first,
		bidExchange:  second,
		askPricePair: exc.PricePair{9600, 100},
		bidPricePair: exc.PricePair{9600, 100},
		orderVolume:  10.0,
	}

	bybitFtx := CrossPair{
		askExchange:  second,
		bidExchange:  first,
		askPricePair: exc.PricePair{9500, 1000},
		bidPricePair: exc.PricePair{9700, 1000},
		orderVolume:  10.0,
	}

	cp := ftxBit
	matchCp := bybitFtx
	cpArray := []CrossPair{
		cp,
	}
	crossPairsTable[cp.GetName()] = cpArray

	matchMap := map[string]CrossPair{}

	matchMap[matchCp.GetName()] = matchCp

	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USDT",
	}

	init := *NewCrossExchangeInit()
	isClose, res := ce.positionCloseCheck(crossPairsTable, matchMap, futures, init)

	assert.True(t, isClose)
	assert.Zero(t, len(res))
}

func TestUserTrades(t *testing.T) {
	symbol := "MAGICUSDT"
	bnfUTs := bnf.GetTightUserTrades(symbol)
	bblUTs := bbl.GetTightUserTrades(symbol)
	specifyTime := time.Now().Add(-time.Hour * 1)
	checkUserTrade(bnfUTs, bblUTs, specifyTime)
	fmt.Println("bbl main")
	checkUserTrade(bblUTs, bnfUTs, specifyTime)

}

func checkUserTrade(mainUTs, otherUTs exc.UserTradeMap, specifyTime time.Time) {
	for key, mainUT := range mainUTs {
		specifyTimeUnixMilli := specifyTime.UnixMilli()
		if key.Time < specifyTimeUnixMilli {
			continue
		}
		if otherUT, ok := otherUTs.Near(key, 1000); ok {

			calc := 1.0
			if mainUT.Side == "BUY" {
				calc = -1.0
			}
			revenue := (mainUT.Price - otherUT.Price) * calc
			revenue = jmath.FloatFloor(revenue, 5)
			fmt.Println(time.UnixMilli(mainUT.Time), revenue, mainUT, otherUT)

		} else {
			fmt.Println(time.UnixMilli(mainUT.Time), "?", mainUT, "-")
		}
	}
}
