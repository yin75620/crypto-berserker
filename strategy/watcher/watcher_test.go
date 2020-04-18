package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func TestWsStart(t *testing.T) {
	iniSetting()
	exchange := ftx.NewFtx(http.DefaultClient,
		ftx.FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			mSubAccount})
	wallet := exchange.GetWallet()
	message := getWalletMessage(wallet)
	fmt.Println(message)
}
