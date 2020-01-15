package CrossExchange

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	simpleLog "github.com/yin75620/crypto-berserker/log"
	"github.com/yin75620/crypto-berserker/message_tool"
)

type CrossExchange struct {
	exchanges        []exc.Exchange
	DelayMilliSecond int64
	futuresArray     []exc.Futures

	//execute use
	positionCrossPairs []CrossPair
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

	message_tool.StartTelegram()

	d := time.Duration(time.Millisecond * time.Duration(ce.DelayMilliSecond))
	//fmt.Println("d1", d)

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C

		plusMilliSecond := ce.stratStrategy()
		//fmt.Println("d1.5", plusMilliSecond)
		t.Reset(time.Millisecond * time.Duration(ce.DelayMilliSecond+plusMilliSecond))
		//fmt.Println("d2", time.Millisecond*time.Duration(ce.DelayMilliSecond+plusMilliSecond))
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

	exchangePairs := []ExchangePair{}

	for _, exchg := range ce.exchanges {
		askPair, bidPair := exchg.GetFuturesAskBidPair(futures)
		exchangePairs = append(exchangePairs, *NewExchangePair(exchg, askPair, bidPair))
	}

	crossPairMap := map[string]CrossPair{}
	for index, ePair := range exchangePairs {
		for matchIndex, matchPair := range exchangePairs {
			if index == matchIndex {
				continue
			}
			cp := NewCrossPair(ePair.exchange, matchPair.exchange, ePair.ask, matchPair.bid)
			crossPairMap[cp.GetName()] = *cp
		}
	}

	// 沒有足夠的交易所可以套利
	if len(crossPairMap) == 0 {
		log.Fatal("has no match pair exchange")
	}

	maxProfit := 0.0
	var topCrossPair CrossPair
	for _, crossPair := range crossPairMap {
		if crossPair.GetProfit() > maxProfit {
			maxProfit = crossPair.GetProfit()
			topCrossPair = crossPair
		}
		crossPair.PrintlnProfit()

		//matchPair := crossPairMap[crossPair.GetMatchName()]
		//totalprofit := matchPair.GetProfit() + crossPair.GetProfit()
	}

	// 沒足夠利潤，直接下一圈
	if maxProfit <= 0 {
		return 0
	}

	// 檢查部位是否可以平倉
	ce.PositionCloseCheck(crossPairMap, futures)

	execMinTotalValue := topCrossPair.GetMinTotalVolume()
	orderTotalValue := math.Min(smallRandom(SETTING_TOTAL_VALUE), execMinTotalValue)
	m_expectedTotalValue = execMinTotalValue - orderTotalValue
	m_expectedLowestProfit = maxProfit

	// 交易判斷
	// 有利可圖
	if !hasProfit(maxProfit, orderTotalValue) {
		// 無利可圖，重設偵測
		m_expectedTotalValue = 0
		m_expectedLowestProfit = 0
		return 0
	}

	// 進行交易
	askExchange, askPair := topCrossPair.GetAskInfo()
	bidExchange, bidPair := topCrossPair.GetBidInfo()

	go executeOrder(askExchange, futures, askPair.Price, exc.Ask, orderTotalValue)
	go executeOrder(bidExchange, futures, bidPair.Price, exc.Bid, orderTotalValue)

	content := fmt.Sprintf("%s, %s\r\n orderTotalValue:%g \r\n maxProfit:%g \r\n m_expectedTotalValue:%g",
		fmt.Sprintf("resAsk:%f, orderVolume:%f, AskCoin:%s", askPair.Price, askPair.Volume, askExchange.GetName()),
		fmt.Sprintf("resBid:%f, orderVolume:%f, bidCoin:%s", bidPair.Price, bidPair.Volume, bidExchange.GetName()),
		orderTotalValue,
		maxProfit,
		m_expectedTotalValue)
	log.Println(content)

	ce.positionCrossPairs = append(ce.positionCrossPairs, topCrossPair)
	message_tool.SendBroadcastArcherGroup(content)

	var plusMilliSecond int64 = 500
	return plusMilliSecond
}

func (ce *CrossExchange) PositionCloseCheck(crossPairMap map[string]CrossPair, futures exc.Futures) {

	//hasPosition
	if len(ce.positionCrossPairs) <= 0 {
		return
	}

	hasRemove := false
	for index, pcp := range ce.positionCrossPairs {
		positionCrossPair := pcp

		matchName := positionCrossPair.GetMatchName()
		if _, ok := crossPairMap[matchName]; !ok {
			log.Println("can't find matchName:", matchName)
			return
		}

		matchCrossPair := crossPairMap[matchName]

		//找出反向配對，確定利潤
		pProfit := positionCrossPair.GetProfit()
		sellProfit := matchCrossPair.GetProfit()
		log.Println(fmt.Sprintf("position profit:%f, sellProfit:%f", pProfit, sellProfit))
		sum := pProfit + sellProfit
		const BASE_PROFIT = 0.0001
		if sellProfit > -0.0007 && sum > BASE_PROFIT {
			log.Println(fmt.Sprintf("sum:%f", sum))

			askExchange, askPair := positionCrossPair.GetAskInfo()
			bidExchange, bidPair := positionCrossPair.GetBidInfo()
			go executeOrder(askExchange, futures, askPair.Price, exc.Ask, askPair.Volume)
			go executeOrder(bidExchange, futures, bidPair.Price, exc.Bid, bidPair.Volume)

			//remove
			hasRemove = true
			ce.positionCrossPairs = removeElement(ce.positionCrossPairs, index)
			break
		}
	}

	// 有資料刪除表示還有資料沒做完，再跑一次。
	if hasRemove {
		ce.PositionCloseCheck(crossPairMap, futures)
	}
}

const (
	SETTING_TOTAL_VALUE = 100
	SETTING_LEVERAGE    = 5.0  //幾倍槓桿
	OverPrice           = 0.03 // 百分之三
)

var (
	m_expectedTotalValue   float64 = 0
	m_expectedLowestProfit float64 = 0
	m_minProfit            float64 = 0.001
	m_minVolume            float64 = 1 //鎂
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

func executeOrder(exchange exc.Exchange, futures exc.Futures, price float64, pType exc.PriceType, volume float64) {
	side := "sell"
	adjPrice := price * (1.0 - OverPrice)
	if pType == exc.Ask {
		side = "buy"
		adjPrice = price * (1.0 + OverPrice)
	}

	myOrder := exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      side,
			Price:     adjPrice,
			Size:      volume,
			OrderType: exc.MARKET,
		},
		Futures: futures,
	}
	exchange.PostFuturesOrder(myOrder)
}

func removeElement(a []CrossPair, i int) []CrossPair {
	copy(a[i:], a[i+1:])      // Shift a[i+1:] left one index.
	a[len(a)-1] = CrossPair{} // Erase last element (write zero value).
	a = a[:len(a)-1]          // Truncate slice.
	return a
}
