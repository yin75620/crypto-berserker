package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-ini/ini"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/binancef"
	"github.com/yin75620/crypto-berserker/exchange-list/bybilinear"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/CrossExchange"
)

const (
	version = "0.9.0-0028-binancef-bybilinear"
)

var (
	mSubAccount string
)

func main() {
	log.Println(version)

	iniSetting()

	exchanges := []exc.Exchange{}
	first := getExchange(common.BINANCEF)
	second := getExchange(common.BYBILINEAR)

	exchanges = append(exchanges, first)
	exchanges = append(exchanges, second)

	ce := CrossExchange.NewCrossExchange(exchanges)
	ce.SetInitByIni("main.ini")
	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USDT",
	}
	ce.SetFuturesArray([]exc.Futures{futures})
	ce.Start()
}

func iniSetting() {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mSubAccount = cfg.Section("FTX").Key("SubAccount").String()
}

func getExchange(exchangeName string) exc.Exchange {
	var res exc.Exchange
	switch exchangeName {
	case common.FTX:
		res = ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			mSubAccount})
		break
	case common.BYBIT:
		res = bybit.NewBybit(http.DefaultClient)
		break
	case common.BYBILINEAR:
		res = bybilinear.NewBybilinear(http.DefaultClient)
		break
	case common.BINANCEF:
		res = binancef.NewBinancef(http.DefaultClient)
		break
	default:
		log.Println("error, not define exchangeName:", exchangeName)
		break
	}

	return res
}
