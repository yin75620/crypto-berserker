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

func (ce *CrossExchange) SetFuturesArray(futuresArray []exc.Futures) {
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

	var cp *CrossPair
	if len(ce.exchanges) < 0 {
		log.Fatal("has no exchange")
	} else if len(ce.exchanges) == 1 {
		cp = NewCrossPair(ce.exchanges[0], ce.exchanges[0])
		log.Println("Only one exchange")
	} else if len(ce.exchanges) > 1 {
		cp = NewCrossPair(ce.exchanges[0], ce.exchanges[1])
	}

	var execProfit, execMinTotalValue float64

	abProfit, abMinTotalValue := cp.GetProfit(futures, MASB)
	execProfit = abProfit
	execMinTotalValue = abMinTotalValue

	baProfit, baMinTotalValue := cp.GetProfit(futures, MBSA)
	if execProfit < baProfit {
		execProfit = baProfit
		execMinTotalValue = baMinTotalValue
	}

	orderTotalValue := math.Min(smallRandom(SETTING_TOTAL_VALUE), execMinTotalValue)
	m_expectedTotalValue = execMinTotalValue - orderTotalValue
	m_expectedLowestProfit = execProfit

	// 交易判斷
	// 有利可圖
	if !hasProfit(execProfit, orderTotalValue) {
		// 無利可圖，重設偵測
		m_expectedTotalValue = 0
		m_expectedLowestProfit = 0
		return 0
	}

	// 進行交易

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

const (
	SETTING_TOTAL_VALUE = 100
)

var (
	m_expectedTotalValue   float64 = 0
	m_expectedLowestProfit float64 = 0
	m_minProfit            float64 = 0.001
	m_minVolume            float64 = 0.00001 //test
	m_tradingAdjustSpeed   int64   = -1000
)

func smallRandom(enter float64) float64 {
	temp := enter * (1 - (10 * rand.Float64() / 100.0)) // 隨機 -10%
	result := math.Floor(temp)
	return result
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
