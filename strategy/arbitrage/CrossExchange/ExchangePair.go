package CrossExchange

import (
	"fmt"
	"log"
	"math"

	exc "github.com/yin75620/crypto-berserker/exchange"
)

type TopPricePair struct {
	ask exc.PricePair
	bid exc.PricePair
}

type ExchangePair struct {
	first TopPricePair
	other TopPricePair
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

type CrossPair struct {
	mainExchange exc.Exchange
	subExchange  exc.Exchange

	crossPosition CrossPosition
}

func NewCrossPair(main, sub exc.Exchange) *CrossPair {
	cp := CrossPair{}
	cp.mainExchange = main
	cp.subExchange = sub
	cp.crossPosition = CrossPosition{}
	return &cp
}

func (cp *CrossPair) GetProfit(futures exc.Futures, ps PositionStatus) (float64, float64) {
	var profit, minSourceTotalValue float64
	switch ps {
	case MASB:
		profit, minSourceTotalValue = cp.abProfit(futures)
	case MBSA:
		profit, minSourceTotalValue = cp.baProfit(futures)
	default:
		log.Println("undefine PositionStatus:", ps)
	}
	return profit, minSourceTotalValue
}

/// 正向 mainAsk(Buy),subBid(Sell)
///profit, minSourceTotalValue
func (cp *CrossPair) abProfit(futures exc.Futures) (float64, float64) {
	mainExchangeName := cp.mainExchange.GetName()
	subExchangeName := cp.subExchange.GetName()

	askPair, _ := cp.mainExchange.GetFuturesAskBidPair(futures)
	askPair.Price = askPair.Price * (1.0 + cp.mainExchange.GetFee().Taker)

	_, bidPair := cp.subExchange.GetFuturesAskBidPair(futures)
	bidPair.Price = bidPair.Price * (1.0 - cp.subExchange.GetFee().Taker)

	log.Println(fmt.Sprintf("ask price:%f, volume:%f, Exchange:%s", askPair.Price, askPair.Volume, mainExchangeName))
	log.Println(fmt.Sprintf("bid price:%f, volume:%f, Exchange:%s", bidPair.Price, bidPair.Volume, subExchangeName))

	return cp.getProfit(askPair, bidPair, mainExchangeName, subExchangeName)
}

/// 反向 mainBid,subAsk
///profit, minSourceTotalValue
func (cp *CrossPair) baProfit(futures exc.Futures) (float64, float64) {
	mainExchangeName := cp.mainExchange.GetName()
	subExchangeName := cp.subExchange.GetName()

	_, bidPair := cp.mainExchange.GetFuturesAskBidPair(futures)
	bidPair.Price = bidPair.Price * (1.0 - cp.mainExchange.GetFee().Taker)

	askPair, _ := cp.subExchange.GetFuturesAskBidPair(futures)
	askPair.Price = askPair.Price * (1.0 + cp.subExchange.GetFee().Taker)

	log.Println(fmt.Sprintf("bid price:%f, volume:%f, Exchange:%s", bidPair.Price, bidPair.Volume, mainExchangeName))
	log.Println(fmt.Sprintf("ask price:%f, volume:%f, Exchange:%s", askPair.Price, askPair.Volume, subExchangeName))

	return cp.getProfit(askPair, bidPair, subExchangeName, mainExchangeName)
}

func (cp *CrossPair) getProfit(askPair, bidPair exc.PricePair, askName, bidName string) (float64, float64) {
	laPrice := askPair.Price
	hbPrice := bidPair.Price

	laVolume := askPair.Volume
	hbVolume := bidPair.Volume
	// 出現錯誤，放慢速度
	if laPrice <= 0 {
		log.Println("laPrice <= 0")
		return 0, 0
	}

	laValue := laPrice * laVolume
	hbValue := hbPrice * hbVolume

	minSourceTotalValue := math.Min(laValue, hbValue)
	log.Println(fmt.Sprintf("minSourceTotalValue:%g", minSourceTotalValue))

	log.Println(fmt.Sprintf("resAsk:%f, laValue:%f, AskName:%s", laPrice, laValue, askName))
	log.Println(fmt.Sprintf("resBid:%f, hbValue:%f, bidName:%s", hbPrice, hbValue, bidName))

	profit := (hbPrice - laPrice) / laPrice

	log.Println(fmt.Sprintf("Profit:%f", profit))

	return profit, minSourceTotalValue
}
