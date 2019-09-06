package main

import (
	"fmt"
	"net/http"

	"github.com/yin75620/crypto-berserker/exchange/ftx"
	fx "github.com/yin75620/crypto-berserker/exchange/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

func main() {
	fmt.Println("TEST")

	var ftx = fx.NewFtx(http.DefaultClient,
		ftx.FtxInit{
			setting.FTX_KEY,
			setting.FTX_API_SECRET_KEY,
			setting.FTX_SUBACCOUNT})

	ftx.GetAccountInfo()

	//ftx.GetMarkets()
}
