package CrossExchange

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	simpleLog "github.com/yin75620/crypto-berserker/log"
)

type CrossExchange struct {
	exchanges        []exc.Exchange
	DelayMilliSecond int64
	futuresArray     []exc.Futures
}

func NewCrossExchange(exchanges []exc.Exchange) *CrossExchange {
	ce := CrossExchange{}
	ce.exchanges = exchanges
	ce.DelayMilliSecond = 500
	return &ce
}

func (ce *CrossExchange) setFuturesArray(futuresArray []exc.Futures) {
	ce.futuresArray = futuresArray
}

func (ce *CrossExchange) Start() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	d := time.Duration(time.Millisecond * time.Duration(ce.DelayMilliSecond))
	fmt.Println("d1", d)

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C

		plusMilliSecond := ce.stratStrategy()
		fmt.Println("d1.5", plusMilliSecond)
		t.Reset(time.Millisecond * time.Duration(ce.DelayMilliSecond+plusMilliSecond))
		fmt.Println("d2", time.Millisecond*time.Duration(ce.DelayMilliSecond+plusMilliSecond))
	}

}

func (ce *CrossExchange) stratStrategy() int64 {

	var totalWaitTime int64
	for _, futures := range ce.futuresArray {
		totalWaitTime += ce.stratFuturesStrategy(futures)
	}
	return totalWaitTime
}

func (ce *CrossExchange) stratFuturesStrategy(futures exc.Futures) int64 {

	laName := ""
	hbName := ""

	//askPairs := []exc.PricePair{}
	//bidPairs := []exc.PricePair{}
	minAskPair := exc.PricePair{Price: math.MaxInt64}
	maxBidPair := exc.PricePair{Price: 0}
	for _, exchg := range ce.exchanges {
		askPair, bidPair := exchg.GetFuturesAskBidPair(futures)
		askPair.Price = askPair.Price * (1.0 + exchg.GetFee().Taker)
		bidPair.Price = bidPair.Price * (1.0 - exchg.GetFee().Taker)

		log.Println(fmt.Sprintf("ask price:%f, volume:%f, Exchange:%s", askPair.Price, askPair.Volume, exchg.GetName()))
		log.Println(fmt.Sprintf("bid price:%f, volume:%f, Exchange:%s", bidPair.Price, bidPair.Volume, exchg.GetName()))

		if askPair.Price < minAskPair.Price {
			minAskPair = askPair
			laName = exchg.GetName()
		}

		if bidPair.Price > maxBidPair.Price {
			maxBidPair = bidPair
			hbName = exchg.GetName()
		}
	}

	// 交易判斷
	ok, delay, orderTotalValue, profit := canOrder(minAskPair, maxBidPair, laName, hbName)
	if !ok {
		return delay
	}

	laOrderVolume := orderTotalValue / minAskPair.Price
	hbOrderVolume := orderTotalValue / maxBidPair.Price

	content := fmt.Sprintf("%s, %s\r\n orderTotalValue:%g \r\n profit:%g \r\n m_expectedTotalValue:%g",
		fmt.Sprintf("resAsk:%f, orderVolume:%f, AskCoin:%s", minAskPair.Price, laOrderVolume, laName),
		fmt.Sprintf("resBid:%f, orderVolume:%f, bidCoin:%s", maxBidPair.Price, hbOrderVolume, hbName),
		orderTotalValue,
		profit,
		m_expectedTotalValue)
	log.Println(content)

	/*const isOrder = true
	if isOrder {
		const isKeepUSD = true
		if isKeepUSD {
			laChannel := tri.executeOrder(lowestAskFlow, exc.Ask, laOrderVolume)
			hbChannel := tri.executeOrder(highestBidFlow, exc.Bid, hbOrderVolume)
			//等上面兩個交易都完成，再繼續
			<-laChannel
			<-hbChannel
		} else {
			hbChannel := tri.executeOrder(highestBidFlow, exc.Bid, laOrderVolume)
			laChannel := tri.executeOrder(lowestAskFlow, exc.Ask, hbOrderVolume)
			//等上面兩個交易都完成，再繼續
			<-laChannel
			<-hbChannel
		}
	}*/

	///

	var plusMilliSecond int64 = 500
	return plusMilliSecond
}

var (
	m_expectedTotalValue   float64 = 0
	m_expectedLowestProfit float64 = 0
	m_minProfit            float64 = 0.001
	m_minVolume            float64 = 20
	m_tradingAdjustSpeed   int64   = -1000
)

func canOrder(lowestAskPair, highestBidPair exc.PricePair, laName, hbName string) (bool, int64, float64, float64) {

	laPrice := lowestAskPair.Price
	hbPrice := highestBidPair.Price

	laVolume := lowestAskPair.Volume
	hbVolume := highestBidPair.Volume
	// 出現錯誤，放慢速度
	if laPrice <= 0 {
		log.Println("laPrice <= 0")
		return false, 60000, 0, 0
	}

	laValue := laPrice * laVolume
	hbValue := hbPrice * hbVolume

	minSourceTotalValue := math.Min(laValue, hbValue)

	log.Println(fmt.Sprintf("minSourceTotalValue:%g", minSourceTotalValue))
	log.Println(fmt.Sprintf("m_expectedTotalValue:%g", m_expectedTotalValue))

	log.Println(fmt.Sprintf("resAsk:%f, laValue:%f, AskCoin:%s", laPrice, laValue, laName))
	log.Println(fmt.Sprintf("resBid:%f, hbValue:%f, bidCoin:%s", hbPrice, hbValue, hbName))

	profit := (hbPrice - laPrice) / laPrice

	log.Println(fmt.Sprintf("Profit:%f", profit))

	currentOrderTotalValue := 100.0 //ce.Init.RANK_N[R_TOTAL_VALUE]

	wnatOrderTotalValue := currentOrderTotalValue * (1 - (10 * rand.Float64() / 100.0)) // 隨機 -10%
	wnatOrderTotalValue = math.Floor(wnatOrderTotalValue)

	orderTotalValue := math.Min(wnatOrderTotalValue, minSourceTotalValue)

	m_expectedTotalValue = minSourceTotalValue - orderTotalValue
	m_expectedLowestProfit = profit

	// 有利可圖
	if !hasProfit(profit, orderTotalValue) {
		// 無利可圖，重設偵測
		m_expectedTotalValue = 0
		m_expectedLowestProfit = 0
		return false, 0, 0, 0
	}

	return true, m_tradingAdjustSpeed, orderTotalValue, profit
}

func hasProfit(profit, orderTotalValue float64) bool {
	// 有利可圖
	if profit < 0 {
		log.Println("No profit")
		return false
	} else if profit < m_minProfit {
		log.Println("No enough profit")
		return false
	} else if orderTotalValue < m_minVolume {
		log.Println(fmt.Sprintf("orderTotalValue < %f", m_minVolume))
		return false
	}
	return true
}

func executeOrder(exchange exc.Exchange, marketName string, price float64, pType exc.PriceType, volume float64) {

	order := exc.ExchangeOrder{}
	order.Market = marketName
	order.Price = price
	order.Size = volume
	side := "sell"
	if pType == exc.Ask {
		side = "buy"
	}
	order.Side = side
	exchange.PostOrder(order)
}
