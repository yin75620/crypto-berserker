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
	"github.com/yin75620/crypto-berserker/setting"
)

type Buyer struct {
	mftx *ftx.Ftx
	init BuyerInit
}

type BuyerInit struct {
	subAccount  string
	targetPrice float64
	size        float64
	coinName    string
}

var (
	mPreWallets []exc.Wallet
)

func main() {

	b := NewBuyer()
	b.SetInit()
	b.InitFtx()
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

	sec := cfg.Section("BUYER")
	init := BuyerInit{}
	init.subAccount = sec.Key("SubAccount").MustString("Buy")
	init.targetPrice = sec.Key("TargetPrice").MustFloat64(1.25)
	init.size = sec.Key("Size").MustFloat64(1)
	init.coinName = sec.Key("CoinName").MustString("MAP")
	b.init = init
}

func (b *Buyer) InitFtx() {
	fi := ftx.NewFtxInit()
	fi.ApiKey = setting.FTX_KEY
	fi.ApiSecretKey = setting.FTX_API_SECRET_KEY
	fi.SubAccount = b.init.subAccount
	b.mftx = ftx.NewFtx(http.DefaultClient, *fi)
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
		CoinPair:  exc.CoinPair{b.init.coinName, "USD"},
		Side:      exc.Buy,
		Price:     b.init.targetPrice,
		Size:      b.init.size,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder)

	return 0
}
