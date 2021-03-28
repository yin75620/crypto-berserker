package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ini/ini"
	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/exchange-list/ftx"
	simpleLog "github.com/yin75620/crypto-berserker/log"
)

const version = "0.9.0-0004"

type Buyer struct {
	mftx *ftx.Ftx
	init BuyerInit
}

type BuyerInit struct {
	subAccount       string
	targetPrice      float64
	size             float64
	coinName         string
	delayMilliSecond float64
	sellPrice        float64
	sellSize         float64
}

var (
	mPreWallets []exc.Wallet
)

func main() {

	b := NewBuyer()
	b.SetInit()
	b.Start()
}

func NewBuyer() *Buyer {
	b := Buyer{}
	return &b
}

func (b *Buyer) SetInit() {
	cfg, err := ini.Load("main.ini")
	if err != nil {
		log.Fatal(err)
	}

	fi := ftx.NewFtxInit()
	fi.IniSettingByFile(cfg)
	b.mftx = ftx.NewFtx(http.DefaultClient, *fi)

	sec := cfg.Section("BUYER")
	init := BuyerInit{}
	init.targetPrice = sec.Key("TargetPrice").MustFloat64(1.25)
	init.size = sec.Key("Size").MustFloat64(1)
	init.coinName = sec.Key("CoinName").MustString("MAP")
	init.delayMilliSecond = sec.Key("DelayMilliSecond").MustFloat64(100)
	init.sellPrice = sec.Key("SellPrice").MustFloat64(3)
	init.sellSize = sec.Key("SellSize").MustFloat64(1)

	b.init = init
}

func (b *Buyer) Start() {

	slog := simpleLog.StartLog()
	defer slog.Close()

	log.Println(version)
	//test api
	infoAll := "Start"

	mftx := b.mftx
	ac := mftx.GetAccountInfo()
	log.Println(string(ac))
	wallet := mftx.GetWallet()
	infoAll = fmt.Sprintf("%s \r\n %s:%v", infoAll, mftx.GetName(), wallet)
	log.Println(string(infoAll))

	delayMilliSecond := time.Millisecond * time.Duration(100)

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

	b.OrderBoth(exchange, "USD")
	b.OrderBoth(exchange, "USDT")

	return 0
}

func (b *Buyer) OrderBoth(exchange *ftx.Ftx, baseCoin string) {
	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{b.init.coinName, baseCoin},
		Side:      exc.Buy,
		Price:     b.init.targetPrice,
		Size:      b.init.size,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder)

	var myOrder2 exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{b.init.coinName, baseCoin},
		Side:      exc.Sell,
		Price:     b.init.sellPrice,
		Size:      b.init.sellSize,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder2)
}
