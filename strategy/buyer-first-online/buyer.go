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

type Buyer struct {
	mftx *ftx.Ftx
}

var (
	mPreWallets []exc.Wallet
)

func main() {
	b := NewBuyer()
	b.Start()
}

func NewBuyer() *Buyer {
	b := Buyer{}

	fi := ftx.NewFtxInit()
	fi.ApiKey = setting.FTX_KEY
	fi.ApiSecretKey = setting.FTX_API_SECRET_KEY
	fi.SubAccount = "Buy"
	b.mftx = ftx.NewFtx(http.DefaultClient, *fi)

	return &b
}

func (b *Buyer) Start() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	//test api
	infoAll := "Start"

	mftx := b.mftx
	ac := mftx.GetAccountInfo()
	log.Println(string(ac))
	wallet := mftx.GetWallet()
	infoAll = fmt.Sprintf("%s \r\n %s:%v", infoAll, mftx.GetName(), wallet)
	log.Println(string(infoAll))

	delayMilliSecond := time.Millisecond * time.Duration(500)

	d := time.Duration(delayMilliSecond)
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C

		b.stratStrategy()
		t.Reset(delayMilliSecond)
	}

}

func (b *Buyer) stratStrategy() int64 {

	exchange := b.mftx

	exchange.DeleteAllOrders()

	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{"MAP", "USD"},
		Side:      exc.Buy,
		Price:     1.25,
		Size:      100,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder)

	return 0
}
