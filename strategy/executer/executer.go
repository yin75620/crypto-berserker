package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yin75620/crypto-berserker/message_tool"

	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/binance"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
	Tri "github.com/yin75620/crypto-berserker/strategy/arbitrage/Triangular"
)

const (
	FTX     = "FTX"
	BITMAX  = "BITMAX"
	MAX     = "MAX"
	OKEX    = "OKEX"
	FTXOTC  = "FTXOTC"
	COINEX  = "COINEX"
	BINANCE = "BINANCE"
)

var mSwitchExchange = BITMAX
var mSubAccount = setting.FTX_SUBACCOUNT

const (
	version = "1.0.1-0013"
)

func main() {

	log.Println(version)

	tInit, bunch := iniSetting()
	s, _ := json.Marshal(tInit)
	log.Println(string(s))

	s2, _ := json.Marshal(bunch)
	log.Println(string(s2))

	log.Println(mSwitchExchange)

	var tri *Tri.Triangular = nil
	var exchange exchange.Exchange = nil

	if mSwitchExchange == FTX {
		exchange = ftx.NewFtx(http.DefaultClient,
			ftx.FtxInit{
				setting.FTX_KEY,
				setting.FTX_API_SECRET_KEY,
				mSubAccount})
	} else if mSwitchExchange == MAX {
		//exchange = maicoin.NewMaicoin(http.DefaultClient)
	} else if mSwitchExchange == OKEX {
		//exchange = okex.NewOkex(http.DefaultClient)
	} else if mSwitchExchange == FTXOTC {
		/*exchange = ftxotc.NewFtxotc(http.DefaultClient,
		ftxotc.FtxotcInit{
			setting.FTXOTC_KEY,
			setting.FTXOTC_SECRET_KEY})*/
	} else if mSwitchExchange == COINEX {
		//exchange = coinex.NewCoinEx(http.DefaultClient)
	} else if mSwitchExchange == BINANCE {
		exchange = binance.NewBinance(http.DefaultClient)

	} else {
	}

	//go dailySendAccountInfo(exchange)

	tri = Tri.NewTriangular(exchange)

	tri.SetInit(tInit)
	tri.SetCoinBunch(bunch)
	tri.Start()

}

func iniSetting() (Tri.TriangularInit, []Tri.CoinStrip) {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mSwitchExchange = cfg.Section("").Key("ExchangeName").String()
	mSubAccount = cfg.Section("FTX").Key("SubAccount").String()

	resInit := Tri.TriangularInit{
		RangePremium:    cfg.Section("").Key("RangePremium").MustFloat64(),
		LeastTotalValue: cfg.Section("").Key("LeastTotalValue").MustFloat64(),
		DelayTime:       cfg.Section("").Key("DelayTime").MustInt(),
		RANK_S: []float64{
			cfg.Section("RankS").Key("MinProfit").MustFloat64(),
			cfg.Section("RankS").Key("PlusSecond").MustFloat64(),
			cfg.Section("RankS").Key("TotalValuesUS").MustFloat64(),
		},
		RANK_N: []float64{
			cfg.Section("RankN").Key("MinProfit").MustFloat64(),
			cfg.Section("RankN").Key("PlusSecond").MustFloat64(),
			cfg.Section("RankN").Key("TotalValuesUS").MustFloat64(),
		},
	}

	resFee := cfg.Section("").Key("TakerFee").MustFloat64()
	log.Println(resFee)

	resFlow := []Tri.CoinStrip{}
	for _, value := range cfg.Section("CoinStrip").Keys() {
		res := value.String()
		log.Println(res)

		coins := strings.Split(res, ",")
		resFlow = append(resFlow, Tri.CoinStrip{Coins: coins})
	}
	return resInit, resFlow
}

// 台灣中午12點會呼叫一次AccountInfo()
func dailySendAccountInfo(exchange exchange.Exchange) {
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	midnoon := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		4, 0, 0, 0, now.Location())
	duration := midnoon.Sub(now)
	time.Sleep(duration)
	content := exchange.GetAccountInfo()
	message_tool.StartTelegram()
	message_tool.SendTelegram(string(content))

	dailySendAccountInfo(exchange)
}
