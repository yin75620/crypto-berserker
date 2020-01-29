package CrossExchange

import (
	"fmt"
	"log"
	"math"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type ExchangePair struct {
	exchange exc.Exchange
	ask      exc.PricePair
	bid      exc.PricePair
}

func NewExchangePair(exchange exc.Exchange, ask, bid exc.PricePair) *ExchangePair {
	ep := ExchangePair{}
	ep.exchange = exchange
	ep.ask = ask
	ep.bid = bid
	return &ep
}

func (ep *ExchangePair) FillPair() {

}

type CrossPair struct {
	askExchange  exc.Exchange
	bidExchange  exc.Exchange
	askPricePair exc.PricePair
	bidPricePair exc.PricePair
	orderVolume  float64
}

func NewCrossPair(ask, bid exc.Exchange, askPircePair, bidPricePair exc.PricePair) *CrossPair {
	cp := CrossPair{}
	cp.askExchange = ask
	cp.bidExchange = bid
	cp.askPricePair = askPircePair
	cp.bidPricePair = bidPricePair
	return &cp
}

func (cp *CrossPair) GetName() string {
	return fmt.Sprintf("A%sB%s", cp.askExchange.GetName(), cp.bidExchange.GetName())
}

func (cp *CrossPair) GetMatchName() string {
	return fmt.Sprintf("A%sB%s", cp.bidExchange.GetName(), cp.askExchange.GetName())
}

func (cp *CrossPair) GetAskInfo() (exc.Exchange, exc.PricePair) {
	return cp.askExchange, cp.askPricePair
}

func (cp *CrossPair) GetBidInfo() (exc.Exchange, exc.PricePair) {
	return cp.bidExchange, cp.bidPricePair
}

//profit
func (cp *CrossPair) GetProfit() float64 {

	aPrice := cp.askPricePair.Price * (1.0 + cp.askExchange.GetFee().Taker)
	bPrice := cp.bidPricePair.Price * (1.0 - cp.bidExchange.GetFee().Taker)

	// 出現錯誤，放慢速度
	if aPrice <= 0 {
		log.Println("laPrice <= 0")
		return 0
	}

	profit := (bPrice - aPrice) / aPrice

	return profit
}

func (cp *CrossPair) GetProfitString() string {

	aPrice := cp.askPricePair.Price * (1.0 + cp.askExchange.GetFee().Taker)
	bPrice := cp.bidPricePair.Price * (1.0 - cp.bidExchange.GetFee().Taker)

	askStr := fmt.Sprintf("ask Cprice:%f, S price:%f, volume:%f, Exchange:%s", aPrice, cp.askPricePair.Price, cp.askPricePair.Volume, cp.askExchange.GetName())
	bidStr := fmt.Sprintf("bid Cprice:%f, S price:%f, volume:%f, Exchange:%s", bPrice, cp.bidPricePair.Price, cp.bidPricePair.Volume, cp.bidExchange.GetName())

	profit := 0.0
	if aPrice > 0 {
		profit = (bPrice - aPrice) / aPrice
	}

	return fmt.Sprintf("A%sB%s Profit:%f \r\n%s, \r\n%s", cp.askExchange.GetName(), cp.bidExchange.GetName(), profit, askStr, bidStr)
}

func (cp *CrossPair) GetMinTotalVolume() float64 {

	aVolume := cp.askPricePair.Volume
	bVolume := cp.bidPricePair.Volume
	// 出現錯誤，放慢速度
	if aVolume <= 0 {
		log.Println("aVolume <= 0")
		return 0
	}

	minTotalVolume := math.Min(aVolume, bVolume)
	log.Println(fmt.Sprintf("minTotalVolume:%g", minTotalVolume))
	return minTotalVolume
}

// M0S0 表示 MainAsk, SubBid 有艙位
// MASB 表示 MainBid, SubAsk 有艙位
// MBSA 表示無艙位
type PositionStatus int

const (
	M0S0 PositionStatus = iota
	MASB
	MBSA
)

type CrossPosition struct {
	mainExchangeName string
	subExchangeName  string
	mainPricePair    exc.PricePair
	subPricePair     exc.PricePair
	status           PositionStatus
}
