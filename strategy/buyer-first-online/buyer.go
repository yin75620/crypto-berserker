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

const version = "0.9.0-0007"

type Buyer struct {
	mftx *ftx.Ftx
	init BuyerInit
}

type BuyOrder struct {
	Price float64
	Size  float64
}

type BuyerInit struct {
	subAccount       string
	targetPrice      float64
	size             float64
	baseCoinName     string
	coinName         string
	delayMilliSecond float64
	sellPrice        float64
	sellSize         float64
	buyOrders        []BuyOrder
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
	init.baseCoinName = sec.Key("BaseCoinName").MustString("USDC")
	init.delayMilliSecond = sec.Key("DelayMilliSecond").MustFloat64(100)
	init.sellPrice = sec.Key("SellPrice").MustFloat64(3)
	init.sellSize = sec.Key("SellSize").MustFloat64(1)

	const DefaultMaxBuy = 9

	orders := []BuyOrder{}
	for i := 0; i < DefaultMaxBuy; i++ {
		sectionBuyStr := fmt.Sprintf("Buy%d", i)
		sectionBuy, err := cfg.GetSection(sectionBuyStr)
		if err != nil {
			continue
		}

		order := BuyOrder{}
		order.Price = sectionBuy.Key("Price").MustFloat64(0.1)
		order.Size = sectionBuy.Key("Size").MustFloat64(1)

		orders = append(orders, order)
	}

	init.buyOrders = orders

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

	delayMilliSecond := time.Millisecond * time.Duration(b.init.delayMilliSecond)

	d := time.Duration(delayMilliSecond)
	t := time.NewTimer(d)
	defer t.Stop()

	exchange := b.mftx

	for {

		exchange.DeleteAllOrders()

		for _, order := range b.init.buyOrders {
			<-t.C
			t.Reset(delayMilliSecond)
			b.OrderPrice(exchange, b.init.baseCoinName, order.Price, order.Size)
		}

		<-t.C
		t.Reset(delayMilliSecond)
		b.SellPrice(exchange, b.init.baseCoinName, b.init.sellPrice, b.init.sellSize)
	}

}

func (b *Buyer) stratStrategy() int64 {

	//	b.OrderBoth(exchange, "USD")
	//b.OrderBoth(exchange, "USDT")

	return 0
}

func (b *Buyer) OrderPrice(exchange *ftx.Ftx, baseCoin string, Price float64, size float64) {
	var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{b.init.coinName, baseCoin},
		Side:      exc.Buy,
		Price:     Price,
		Size:      size,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder)
}

func (b *Buyer) SellPrice(exchange *ftx.Ftx, baseCoin string, Price float64, size float64) {
	var myOrder2 exc.ExchangeOrder = exc.ExchangeOrder{
		CoinPair:  exc.CoinPair{b.init.coinName, baseCoin},
		Side:      exc.Sell,
		Price:     Price,
		Size:      size,
		OrderType: exc.LIMIT,
	}
	exchange.PostOrder(myOrder2)
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
