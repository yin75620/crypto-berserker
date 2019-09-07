package exchange

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Exchange interface {
	GetAskBidPair(coinPair CoinPair, depth int) (PricePair, PricePair)
	GetAccountInfo() []byte
	PostOrder(order ExchangeOrder) (string, error)
	GetFee() Fee
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

type EOrderType string

const (
	LIMIT  EOrderType = "limit"
	MARKET EOrderType = "market"
)

type ExchangeOrder struct {
	CoinPair  CoinPair
	Market    string     `json:"market"`
	Side      string     `json:"side"`
	Price     float64    `json:"price"`
	Size      float64    `json:"size"`
	OrderType EOrderType `json:"order_type"`
	//ReduceOnly bool       `json:"reduceOnly"`
}

type CoinPair struct {
	BaseCoin   string //基礎貨幣
	QuotedCoin string //標價貨幣
}

func (co *CoinPair) GetMarketName() string {
	return fmt.Sprintf("%s/%s", co.BaseCoin, co.QuotedCoin)
}

func (co *CoinPair) GetSymbal() string {
	return fmt.Sprintf("%s-%s", co.BaseCoin, co.QuotedCoin)
}

func (co *CoinPair) GetLinkMakertName() string {
	return strings.ToLower(fmt.Sprintf("%s%s", co.BaseCoin, co.QuotedCoin))
}

func SendRequest(client *http.Client, req *http.Request) []byte {
	var res []byte
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return res
	}

	defer resp.Body.Close()
	sitemap, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("%s", err)
		return res
	}

	log.Printf("%s", sitemap)
	return sitemap
}
