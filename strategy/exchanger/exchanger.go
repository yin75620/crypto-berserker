package main

import (
	"log"
	"net/http"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/bybit"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	"github.com/yin75620/crypto-berserker/setting"
	"github.com/yin75620/crypto-berserker/strategy/arbitrage/CrossExchange"
)

const (
	version = "0.9.0-0008"
)

func main() {
	log.Println(version)

	exchanges := []exc.Exchange{}
	ft := ftx.NewFtx(http.DefaultClient, ftx.FtxInit{
		setting.FTX_KEY,
		setting.FTX_API_SECRET_KEY,
		"Archer"})
	bybit := bybit.NewBybit(http.DefaultClient)
	exchanges = append(exchanges, ft)
	exchanges = append(exchanges, bybit)

	ce := CrossExchange.NewCrossExchange(exchanges)
	futures := exc.Futures{
		//ExpirationDate: time.Date(2019, time.December, 27, 0, 0, 0, 0, time.UTC),
		TargetName: "BTC",
		QuoteCoin:  "USD",
	}
	ce.SetFuturesArray([]exc.Futures{futures})
	ce.Start()
}
