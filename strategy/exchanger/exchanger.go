package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/common"
	"github.com/yin75620/crypto-berserker/message_tool"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/CrossExchange"
)

const (
	version = "0.9.2-0043"
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

	sendWallet(exchanges)
	go dailySendAccountInfo(exchanges)

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

// 台灣中午12點會呼叫一次AccountInfo()
func dailySendAccountInfo(exchanges []exc.Exchange) {
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	midnoon := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		4, 0, 0, 0, now.Location())
	duration := midnoon.Sub(now)
	time.Sleep(duration)
	sendWallet(exchanges)

	dailySendAccountInfo(exchanges)
}

func sendWallet(exchanges []exc.Exchange) {
	message := ""
	sumUSD := 0.0
	for _, v := range exchanges {
		wallet := v.GetWallet()
		sumUSD += wallet.GetAllBalanceUSDValue()
		message += fmt.Sprintf("%s \r\n", v.GetName())
		message += wallet.GetWalletMessage()
		message += "\r\n"
	}
	message += fmt.Sprintf("sumUSD:%.2f", sumUSD)
	message_tool.StartTelegram()
	message_tool.SendWatcherGroup(message)
}
