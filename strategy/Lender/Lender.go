package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yin75620/crypto-berserker/message_tool"

	"github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
)

const version = "0.9.0-0003"

const IniFileName = "main.ini"

var mSubAccount = ""

func main() {

	log.Println(version)

	fi := ftx.NewFtxInit()
	fi.IniSetting(IniFileName)
	mSubAccount = fi.SubAccount
	myFtx := ftx.NewFtx(http.DefaultClient, *fi)

	// get first on the our + 1 minute
	const HOUR = 60
	waitMinute := HOUR - time.Now().Minute() + 1
	hour := time.Minute * time.Duration(waitMinute)
	hTimer := time.NewTimer(hour)
	defer hTimer.Stop()

	min := time.Minute * time.Duration(1)
	minTimer := time.NewTimer(min)
	defer minTimer.Stop()

	doOffer(*myFtx)

	for {
		select {
		case <-hTimer.C:
			doOffer(*myFtx)

			hTimer.Reset(time.Hour * time.Duration(1))
		case <-minTimer.C:
			log.Println("Live")
		}
	}
}

func doOffer(myFtx ftx.Ftx) {
	const coin = "USD"
	res := myFtx.GetLendingInfo()
	li := res.GetLendInfo("USD")
	fmt.Println(res)

	if li.Offerd != 0 {
		lo := ftx.NewLendOrder("USD", li.Lendable, li.MinRate)
		myFtx.PostLendingOffer(*lo)
	}

	res = myFtx.GetLendingInfo()
	fmt.Println(res)
}

// 台灣中午12點會呼叫一次AccountInfo()
func dailySendAccountInfo(exchange exchange.Exchange) {
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	midnoon := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		4, 0, 0, 0, now.Location())
	duration := midnoon.Sub(now)
	time.Sleep(duration)
	sendWallet(exchange)

	dailySendAccountInfo(exchange)
}

func sendWallet(exchange exchange.Exchange) {
	wallet := exchange.GetWallet()
	message := fmt.Sprintf("%s (%s) \r\n", exchange.GetName(), mSubAccount)
	message += wallet.GetWalletMessage()
	message_tool.StartTelegram()
	message_tool.SendWatcherGroup(message)
}
