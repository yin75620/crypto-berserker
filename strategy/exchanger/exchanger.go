package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-ini/ini"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/CrossExchange"
)

const (
	version = "0.9.1-0035"
)

var (
	mExchangeStrings []string
	mFutures         exc.Futures
)

func main() {
	log.Println(version)

	iniSetting()

	exchanges := []exc.Exchange{}

	for _, v := range mExchangeStrings {
		ex := common.GetExchange(v)
		exchanges = append(exchanges, ex)
	}

	ce := CrossExchange.NewCrossExchange(exchanges)
	ce.SetInitByIni("main.ini")
	ce.SetFuturesArray([]exc.Futures{mFutures})
	ce.Start()
}

func iniSetting() {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	res := cfg.Section("Exchanger").Key("Exchanges").String()
	mExchangeStrings = strings.Split(res, ",")

	targetName := cfg.Section("Futures").Key("TargetName").String()
	quoteCoin := cfg.Section("Futures").Key("QuoteCoin").String()

	mFutures = exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: targetName,
		QuoteCoin:  quoteCoin,
	}
}
