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
	version = "1.0.1-0016"
)

func main() {

	log.Println(version)

	tInit, bunchs := iniSetting()
	s, _ := json.Marshal(tInit)
	log.Println(string(s))

	s2, _ := json.Marshal(bunchs)
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
	tri.SetCoinBunchs(bunchs)
	tri.Start()

}

func iniSetting() (Tri.TriangularInit, []Tri.CoinBunch) {
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
	}

	resFee := cfg.Section("").Key("TakerFee").MustFloat64()
	log.Println(resFee)

	const DefaultMaxCoinStrip = 10
	resFlow := []Tri.CoinBunch{}

	for i := 0; i < DefaultMaxCoinStrip; i++ {
		coinBunch := Tri.CoinBunch{}
		coinStrips := []Tri.CoinStrip{}
		sectionStripStr := fmt.Sprintf("CoinStrip%d", i)
		sectionStrip, err := cfg.GetSection(sectionStripStr)
		if err != nil {
			continue
		}
		for _, value := range sectionStrip.Keys() {
			if !strings.HasPrefix(value.Name(), "key") {
				continue
			}
			res := value.String()
			log.Println(res)

			coins := strings.Split(res, ",")
			coinStrips = append(coinStrips, Tri.CoinStrip{Coins: coins})
		}
		coinBunch.CoinStrips = coinStrips
		coinBunch.MinProfit = sectionStrip.Key("MinProfit").MustFloat64()
		coinBunch.PlusSecond = sectionStrip.Key("PlusSecond").MustFloat64()
		coinBunch.TotalValuesUS = sectionStrip.Key("TotalValuesUS").MustFloat64()
		resFlow = append(resFlow, Tri.CoinBunch{CoinStrips: coinStrips})
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
