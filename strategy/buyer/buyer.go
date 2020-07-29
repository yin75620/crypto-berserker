package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	simpleLog "github.com/yin75620/crypto-berserker/log"
	"github.com/yin75620/crypto-berserker/setting"
)

type CrossExchange struct {
	exchange exc.Exchange
}

var (
	mPreWallets []exc.Wallet
)

func main() {
	ce := NewCrossExchange()
	ce.Start()
}

func NewCrossExchange() *CrossExchange {
	ce := CrossExchange{}

	fi := ftx.NewFtxInit()
	fi.ApiKey = setting.FTX_KEY
	fi.ApiSecretKey = setting.FTX_API_SECRET_KEY
	fi.SubAccount = "DeFi Competition"
	exchange := ftx.NewFtx(http.DefaultClient, *fi)

	ce.exchange = exchange
	return &ce
}

func (ce *CrossExchange) Start() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	//test api
	infoAll := "Start"

	exchg := ce.exchange
	ac := exchg.GetAccountInfo()
	log.Println(string(ac))
	wallet := exchg.GetWallet()
	infoAll = fmt.Sprintf("%s \r\n %s:%v", infoAll, exchg.GetName(), wallet)
	log.Println(string(infoAll))

	delayMilliSecond := time.Millisecond * time.Duration(5000)

	d := time.Duration(delayMilliSecond)
	t := time.NewTimer(d)
	defer t.Stop()

	//startTime := time.Now()

	for {
		<-t.C

		ce.stratStrategy()
		t.Reset(delayMilliSecond)
	}

}

func (ce *CrossExchange) stratStrategy() int64 {

	exchange := ce.exchange

	futures := exc.Futures{}
	futures.QuoteCoin = "USD"
	futures.TargetName = "CUSDT"
	ap, bp := exchange.GetFuturesAskBidPair(futures)

	midPrice := (ap.Price + bp.Price) * 0.5

	tick1 := 0.0000025

	plus1Price := midPrice + tick1
	minus1Price := midPrice - tick1

	ce.postOrder(exchange, "buy", 0.020220, 10000, futures)
	ce.postOrder(exchange, "sell", 0.0202050, 10000, futures)

	return 0
}

func (ce *CrossExchange) postOrder(exchange exc.Exchange, side string, adjPrice, volume float64, futures exc.Futures) {

	myOrder := exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      side,
			Price:     adjPrice,
			Size:      volume,
			OrderType: exc.LIMIT,
		},
		Futures: futures,
		IsClose: false,
	}
	exchange.PostFuturesOrder(myOrder)

}
