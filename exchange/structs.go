package exchange

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/yin75620/crypto-berserker/exchange/tool"
)

type JArray map[string]interface{}

func (ja *JArray) Add(other JArray) {
	// 塞入 body
	for k, v := range other {
		(*ja)[k] = v
	}
}

func (ja *JArray) ToValues() url.Values {
	// 塞入 values
	values := url.Values{}
	for k, v := range *ja {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return values
}

type ErrorMessage struct {
	Code    int    `"json:code,omitempty"`
	Message string `"json:message,omitempty"`
}

type PriceStatus struct {
	Asks [][]float64 `"json:asks,string,omitempty"`
	Bids [][]float64 `"json:bids,string,omitempty"`
}

func (ps *PriceStatus) GetPair(depth int, pType PriceType) (PricePair, error) {
	switch pType {
	case Ask:
		return ps.getAskPricePair(depth)
	case Bid:
		return ps.getBidPricePair(depth)
	}
	return PricePair{}, errors.New("has no match PriceType")
}

func (ps *PriceStatus) getAskPricePair(depth int) (PricePair, error) {
	return GetPricePair(depth, ps.Asks)
}
func (ps *PriceStatus) getBidPricePair(depth int) (PricePair, error) {
	return GetPricePair(depth, ps.Bids)
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

func (ps *PriceStatus) SetByJArray(json map[string]interface{}) {
	ps.Asks = tool.TransToFloatTwoArray(json["asks"].([]interface{}))
	ps.Bids = tool.TransToFloatTwoArray(json["bids"].([]interface{}))
}

type ICommodity interface {
	GetMarketName()
}

type Spot struct {
	//ICommodity
	CoinPair
}

type Futures struct {
	//ICommodity
	//到期日
	ExpirationDate time.Time
	// 商品名
	TargetName string
	// 計價貨幣類型
	QuoteCoin string
}

func (f *Futures) GetLinkMarketName() string {
	res := ""
	if f.ExpirationDate.IsZero() {
		res = fmt.Sprintf("%s%s", f.TargetName, f.QuoteCoin)
	}
	return res
}

func (f *Futures) GetMarketName() string {
	res := ""
	if f.ExpirationDate.IsZero() {
		res = fmt.Sprintf("%s-PERP", f.TargetName)
	} else {
		res = fmt.Sprintf("%s-%d%d", f.TargetName, f.ExpirationDate.Month(), f.ExpirationDate.Day())
	}

	return res
}

func (f *Futures) GetSwapNameUpper() string {
	return strings.ToUpper(fmt.Sprintf("%s-%s-SWAP", f.TargetName, f.QuoteCoin))
}
