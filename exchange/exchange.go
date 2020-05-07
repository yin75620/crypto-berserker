package exchange

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
)

type Exchange interface {
	GetAccountInfo() []byte // include balance
	PostOrder(order ExchangeOrder) (string, error)
	PostFuturesOrder(order FuturesOrder) (string, error)
	GetWallet() Wallet
	GetFee() Fee
	GetName() string
	GetMarketInfo(coinPair CoinPair) MarketInfo
	GetAskBidPair(coinPair CoinPair, depth int) (PricePair, PricePair)
	GetFuturesAskBidPair(futures Futures) (PricePair, PricePair) // VOLUME = USD
	GetVolumeByTotal(total, price float64) float64               //volume

	GetAccount() Account
}

func PostOrderRefry(ex Exchange, order ExchangeOrder) {
	//最多重試三次
	const RETRY_TIMES = 3
	for i := 0; i < RETRY_TIMES; i++ {
		_, err := ex.PostOrder(order)
		if err == nil {
			break
		}
		log.Println(fmt.Sprintf("retry Count:%d", i))
	}
}

func PostFuturesOrderRefry(ex Exchange, order FuturesOrder) {
	//最多重試三次
	const RETRY_TIMES = 3
	for i := 0; i < RETRY_TIMES; i++ {
		_, err := ex.PostFuturesOrder(order)
		if err == nil {
			break
		}
		log.Println(fmt.Sprintf("retry Count:%d", i))
	}
}

type PriceType int

const (
	Ask PriceType = iota
	Bid
)

const (
	Sell = "sell"
	Buy  = "buy"
)

type PricePair struct {
	Price  float64
	Volume float64
}

func (p *PricePair) Total() float64 {
	return p.Price * p.Volume
}

func (p *PricePair) CalcuVolumeToTatol() {
	p.Volume = p.Volume / p.Price
}

func LowestPricePair(pps []PricePair) PricePair {
	minAskPair := PricePair{}
	minAskPair.Price = math.MaxInt64
	for _, ask := range pps {
		if ask.Price < minAskPair.Price {
			minAskPair = ask
		}
	}
	return minAskPair
}

func HighestPricePair(pps []PricePair) PricePair {
	bidPair := PricePair{}
	bidPair.Price = 0
	for _, ask := range pps {
		if ask.Price > bidPair.Price {
			bidPair = ask
		}
	}
	return bidPair
}

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type ExchangeOrder struct {
	CoinPair CoinPair // will drop this parameter
	Market   string   `json:"market"` //temp
	// 預計把上面兩個換成 Commodity
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

type FuturesOrder struct {
	CommodityOrder
	Futures Futures
}

type CommodityOrder struct {
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
}

type CoinPair struct {
	BaseCoin   string //基礎貨幣
	QuotedCoin string //標價貨幣
}

func (co *CoinPair) GetMarketName() string {
	return fmt.Sprintf("%s/%s", co.BaseCoin, co.QuotedCoin)
}

func (co *CoinPair) SetByMarketName(marketName string) {
	res := strings.Split(marketName, "/")
	co.BaseCoin = res[0]
	co.QuotedCoin = res[1]
}

func (co *CoinPair) GetSymbal() string {
	return fmt.Sprintf("%s-%s", co.BaseCoin, co.QuotedCoin)
}

func (co *CoinPair) SetBySymbal(symbal string) {
	res := strings.Split(symbal, "-")
	co.BaseCoin = res[0]
	co.QuotedCoin = res[1]
}

func (co *CoinPair) GetLinkMakertName() string {
	return strings.ToLower(fmt.Sprintf("%s%s", co.BaseCoin, co.QuotedCoin))
}

func (co *CoinPair) GetLinkMakertNameUpper() string {
	return strings.ToUpper(fmt.Sprintf("%s%s", co.BaseCoin, co.QuotedCoin))
}

func (co *CoinPair) SetByLinkMakertName(linkName string) {
	quoteCoins := []string{"btc", "usdt", "twd", "usd"}
	for _, v := range quoteCoins {
		index := strings.LastIndex(linkName, v)
		if index == -1 {
			continue
		}

		length := len(linkName)
		co.BaseCoin = strings.ToUpper(linkName[0:index])
		co.QuotedCoin = strings.ToUpper(linkName[index:length])
		return
	}
	fmt.Println("new linkName:", linkName)
}

func SendRequest(client *http.Client, req *http.Request) ([]byte, error) {
	var res []byte
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res, err
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("%s", err)
		return res, err
	}

	//log.Printf("%s", sitemap)
	return sitemap, nil
}

type MarketInfo struct {
	Name            string
	PriceIncrement  float64
	VolumeIncrement float64
}
