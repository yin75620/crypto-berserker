package exchange

import (
	"errors"
	"fmt"
)

type Exchange interface {
	GetAskBidPair(coinPair CoinPair, depth int) (PricePair, PricePair)
	GetAccountInfo() []byte
	PostOrder(order ExchangeOrder) string
}

type PriceType int

const (
	Ask PriceType = iota
	Bid
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

type PriceStatus interface {
	GetPair(depth int, pType PriceType) (PricePair, error)
	getAskPricePair(depth int) (PricePair, error)
	getBidPricePair(depth int) (PricePair, error)
}

func GetPricePair(depth int, prices [][]float64) (PricePair, error) {
	var res = PricePair{}
	size := len(prices)
	if depth > size {
		return res, errors.New("depth can't over size")
	}

	index := depth - 1
	res.Price = prices[index][0] // first prize, second volume
	res.Volume = prices[index][1]
	return res, nil
}
