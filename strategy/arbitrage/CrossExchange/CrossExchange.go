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
	//positionCrossPairs []CrossPair
	positionCrossPairMap map[string][]CrossPair
}

func NewCrossExchange(exchanges []exc.Exchange) *CrossExchange {
	ce := CrossExchange{}
	ce.exchanges = exchanges
	ce.DelayMilliSecond = 500
	ce.positionCrossPairMap = map[string][]CrossPair{}
	return &ce
}

func (ce *CrossExchange) SetFuturesArray(futuresArray []exc.Futures) {
	ce.futuresArray = futuresArray
}

func (ce *CrossExchange) Start() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	message_tool.StartTelegram()

	//test api
	infoAll := ""
	for _, exchg := range ce.exchanges {
		info := exchg.GetAccountInfo()
		infoAll = fmt.Sprintf("%s \r\n %s", infoAll, info)
	}
	log.Println(string(infoAll))
	message_tool.SendBroadcastArcherGroup(infoAll)

	d := time.Duration(time.Millisecond * time.Duration(ce.DelayMilliSecond))

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

	crossPairMap := getCrossPairMap(ce.exchanges, futures)

	// 沒有足夠的交易所可以套利
	if len(crossPairMap) == 0 {
		log.Fatal("has no match pair exchange")
	}

	topCrossPair := getMaxProfitCrossPair(crossPairMap)
	maxProfit := topCrossPair.GetProfit()

	// 檢查部位是否可以平倉
	positionCloseCheck(ce.positionCrossPairMap, crossPairMap, futures)

	execMinTotalValue := topCrossPair.GetMinTotalVolume()
	orderTotalValue := math.Min(MAX_HOLD_VOLUME-m_currentUSDVolume, execMinTotalValue)

	// 交易判斷
	// 有利可圖
	if !canOrder(maxProfit, orderTotalValue) {
		// 無利可圖，重設偵測
		return 0
	}

	m_expectedTotalValue = execMinTotalValue - orderTotalValue
	m_expectedLowestProfit = maxProfit

	// 進行交易
	orderCrossPair(topCrossPair, futures, orderTotalValue)
	ce.positionCrossPairMap[topCrossPair.GetName()] = append(ce.positionCrossPairMap[topCrossPair.GetName()], topCrossPair)

	var plusMilliSecond int64 = 500
	return plusMilliSecond
}

func positionCloseCheck(crossPairsTable map[string][]CrossPair, matchMap map[string]CrossPair, futures exc.Futures) map[string][]CrossPair {

	//hasPosition
	if len(crossPairsTable) <= 0 {
		return crossPairsTable
	}

	for key, arrayPairs := range crossPairsTable {

		closeCheck(key, arrayPairs, matchMap, futures)

		for i := len(arrayPairs) - 1; i >= 0; i-- {
			pair := arrayPairs[i]
			if pair.orderVolume == 0 {
				//remove
				arrayPairs = removeElement(arrayPairs, i)
			}
		}

		if len(arrayPairs) == 0 {
			delete(crossPairsTable, key)
		} else {
			crossPairsTable[key] = arrayPairs
		}
	}

	return crossPairsTable
}

func getMatchCrossPair(positionPairName string, crossPairArray []CrossPair, matchMap map[string]CrossPair) CrossPair {
	// key check
	matchName := ""
	for _, crossPair := range crossPairArray {
		if positionPairName != crossPair.GetName() {
			log.Println(fmt.Sprintf("map has not same key:%s, positionPairName:%s", crossPair.GetName(), positionPairName))
			return CrossPair{}
		}

		// MatchName 都一樣，選一個 matchName 放入
		matchName = crossPair.GetMatchName()
	}

	if len(crossPairArray) <= 0 {
		log.Println("closeCheck len(crossPairArray) <= 0 ")
		return CrossPair{}
	}

	matchCrossPair, ok := matchMap[matchName]
	if !ok {
		log.Fatalf("matchMap not found. matchName:%s", matchName)
		return CrossPair{}
	}
	return matchCrossPair
}

func getTotalVolume(crossPairArray []CrossPair, matchCrossPair CrossPair) (float64, float64, float64) {
	matchProfit := matchCrossPair.GetProfit()
	matchVolume := matchCrossPair.GetMinTotalVolume()
	// 確認利潤並統計總量
	askTotalVolume := 0.0
	bidTotalVolume := 0.0
	totalMatchOrderUSDVolume := 0.0
	for _, pcp := range crossPairArray {
		positionCrossPair := pcp

		//找出反向配對，確定利潤
		positionProfit := positionCrossPair.GetProfit()

		log.Println(fmt.Sprintf("position profit:%f, matchProfit:%f", positionProfit, matchProfit))
		sum := positionProfit + matchProfit
		if matchProfit > MinSellProfit && sum > MinSumProfit && matchVolume > m_minVolume {
			log.Println(fmt.Sprintf("sum:%f", sum))

			thisMatchOrderVolume := math.Min(matchVolume, positionCrossPair.orderVolume)

			askVolume := positionCrossPair.GetAskVolumeByTotal(thisMatchOrderVolume)
			bidVolume := positionCrossPair.GetBidVolumeByTotal(thisMatchOrderVolume)

			askTotalVolume += askVolume
			bidTotalVolume += bidVolume

			totalMatchOrderUSDVolume += thisMatchOrderVolume

			content := fmt.Sprintf(
				"positionCrossPair:%s,\r\n matchPair:%s\r\n sumProfit:%f",
				positionCrossPair.GetProfitString(),
				matchCrossPair.GetProfitString(),
				sum)
			message_tool.SendBroadcastArcherGroup(content)

			positionCrossPair.orderVolume -= thisMatchOrderVolume
		}
	}
	return askTotalVolume, bidTotalVolume, totalMatchOrderUSDVolume
}

func closeCheck(positionPairName string, crossPairArray []CrossPair, matchMap map[string]CrossPair, futures exc.Futures) {

	matchCrossPair := getMatchCrossPair(positionPairName, crossPairArray, matchMap)

	askTotalVolume, bidTotalVolume, totalMatchOrderUSDVolume := getTotalVolume(crossPairArray, matchCrossPair)

	// 表示有交易
	if totalMatchOrderUSDVolume > 0 {
		matchAskExchange, askPair := matchCrossPair.GetAskInfo()
		matchBidExchange, bidPair := matchCrossPair.GetBidInfo()

		askChannel := executeOrder(matchAskExchange, futures, askPair.Price, exc.Ask, bidTotalVolume)
		bidChannel := executeOrder(matchBidExchange, futures, bidPair.Price, exc.Bid, askTotalVolume)
		//等上面兩個交易都完成，再繼續
		<-askChannel
		<-bidChannel

		m_currentUSDVolume = m_currentUSDVolume - totalMatchOrderUSDVolume
	}

}

func getCrossPairMap(exchanges []exc.Exchange, futures exc.Futures) map[string]CrossPair {

	exchangePairs := []ExchangePair{}
	for _, exchg := range exchanges {
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
	return crossPairMap
}

func getMaxProfitCrossPair(mcp map[string]CrossPair) CrossPair {
	maxProfit := -math.MaxFloat64
	var topCrossPair CrossPair
	for _, crossPair := range mcp {
		if crossPair.GetProfit() > maxProfit {
			maxProfit = crossPair.GetProfit()
			topCrossPair = crossPair
		}
		log.Println(crossPair.GetProfitString())

		//matchPair := crossPairMap[crossPair.GetMatchName()]
		//totalprofit := matchPair.GetProfit() + crossPair.GetProfit()
	}
	return topCrossPair
}

func orderCrossPair(topCrossPair CrossPair, futures exc.Futures, orderTotalValue float64) {
	askExchange, askPair := topCrossPair.GetAskInfo()
	bidExchange, bidPair := topCrossPair.GetBidInfo()

	askVolume := askExchange.GetVolumeByTotal(orderTotalValue, askPair.Price)
	bidVolume := bidExchange.GetVolumeByTotal(orderTotalValue, bidPair.Price)

	askChannel := executeOrder(askExchange, futures, askPair.Price, exc.Ask, askVolume)
	bidChannel := executeOrder(bidExchange, futures, bidPair.Price, exc.Bid, bidVolume)
	//等上面兩個交易都完成，再繼續
	<-askChannel
	<-bidChannel

	m_currentUSDVolume = m_currentUSDVolume + orderTotalValue

	content := fmt.Sprintf("%s, %s\r\n orderTotalValue:%g \r\n maxProfit:%g \r\n m_expectedTotalValue:%g",
		fmt.Sprintf("resAsk:%f, orderVolume:%f, AskCoin:%s", askPair.Price, askPair.Volume, askExchange.GetName()),
		fmt.Sprintf("resBid:%f, orderVolume:%f, bidCoin:%s", bidPair.Price, bidPair.Volume, bidExchange.GetName()),
		orderTotalValue,
		topCrossPair.GetProfit(),
		m_expectedTotalValue)
	log.Println(content)

	// 調整成交量，改用下單的量，後續平倉成交量才會正確。
	topCrossPair.orderVolume = orderTotalValue

	message_tool.SendBroadcastArcherGroup(content)
}

const (
	//SETTING_TOTAL_VALUE = 1510.0
	//SETTING_LEVERAGE    = 5.0  //幾倍槓桿
	OverPrice     = 0.02 // 交易時，要溢價多少。 Ex:目前價位 9000 => 會用9180買進
	MinSellProfit = -0.0007
	MinSumProfit  = 0.0001

	MAX_HOLD_VOLUME = 1000.0
)

var (
	m_expectedTotalValue   float64 = 0
	m_expectedLowestProfit float64 = 0
	m_minProfit            float64 = 0.001
	m_minVolume            float64 = 1 //鎂
	m_tradingAdjustSpeed   int64   = -1000

	m_currentUSDVolume float64 = 0

	m_cutPercent float64 = 10 //隨機減少N%下單N>0
)

func smallRandom(enter float64) float64 {
	temp := enter * (1 - (m_cutPercent * rand.Float64() / 100.0)) // 隨機 -10%
	result := math.Floor(temp)
	return result
}

func canOrder(profit, orderTotalValue float64) bool {
	// 有利可圖
	if profit < m_minProfit { // 沒足夠利潤，直接下一圈
		log.Println(fmt.Sprintf("No enough profit. profit:%f", profit))
		return false
	} else if orderTotalValue < m_minVolume {
		log.Println(fmt.Sprintf("orderTotalValue < %f", m_minVolume))
		return false
	} else if m_currentUSDVolume >= MAX_HOLD_VOLUME {
		log.Println(fmt.Sprintf("m_currentUSDVolume:%g >= MAX_HOLD_VOLUME:%g", m_currentUSDVolume, MAX_HOLD_VOLUME))
		return false
	}
	return true
}

func executeOrder(exchange exc.Exchange, futures exc.Futures, price float64, pType exc.PriceType, volume float64) chan int {
	resultChannel := make(chan int)
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
	go func() {
		exchange.PostFuturesOrder(myOrder)
		resultChannel <- 0
	}()
	return resultChannel
}

func removeElement(a []CrossPair, i int) []CrossPair {
	copy(a[i:], a[i+1:])      // Shift a[i+1:] left one index.
	a[len(a)-1] = CrossPair{} // Erase last element (write zero value).
	a = a[:len(a)-1]          // Truncate slice.
	return a
}
