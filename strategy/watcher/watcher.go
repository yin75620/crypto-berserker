package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yin75620/crypto-berserker/message_tool"

	"github.com/go-ini/ini"
	"github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/binance"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
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
	version = "0.9.0-0001"
)

func main() {

	log.Println(version)

	iniSetting()

	log.Println(mSwitchExchange)

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

	dailySendAccountInfo(exchange)

}

func iniSetting() {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mSwitchExchange = cfg.Section("").Key("ExchangeName").String()
	mSubAccount = cfg.Section("FTX").Key("SubAccount").String()
}

// 台灣中午12點會呼叫一次AccountInfo()
func dailySendAccountInfo(exchange exchange.Exchange) {
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	midnoon := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		4, 0, 0, 0, now.Location())
	duration := midnoon.Sub(now)
	time.Sleep(duration)
	wallet := exchange.GetWallet()
	message := getWalletMessage(wallet)
	message_tool.StartTelegram()
	message_tool.SendWatcherGroup(message)

	dailySendAccountInfo(exchange)
}

func getWalletMessage(wallet exchange.Wallet) string {
	message := ""
	for _, v := range wallet.Balances {
		message += fmt.Sprintf("%4s, %.2f, US$%.2f\r\n", v.Coin, v.Total, v.UsdValue)
	}
	message += fmt.Sprintf("Total:US$%.2f", wallet.GetAllBalanceUSDValue())
	return message
}
