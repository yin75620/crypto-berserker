package CrossExchange

import (
	"fmt"
	"log"
	"math"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	simpleLog "github.com/yin75620/crypto-berserker/log"
	"github.com/yin75620/crypto-berserker/message_tool"
)

type CrossPairStringMap map[string]CrossPair

type CrossExchange struct {
	exchanges            []exc.Exchange
	futuresArray         []exc.Futures
	positionCrossPairMap map[string][]CrossPair
	init                 CrossExchangeInit
}

func NewCrossExchange(exchanges []exc.Exchange) *CrossExchange {
	ce := CrossExchange{}
	ce.exchanges = exchanges
	ce.positionCrossPairMap = map[string][]CrossPair{}
	NewCrossExchangeInit()
	ce.init = *NewCrossExchangeInit()
	return &ce
}

func (ce *CrossExchange) SetInit(init CrossExchangeInit) {
	ce.init = init
}

func (ce *CrossExchange) SetInitByIni(filename string) {
	ce.init.IniSetting(filename)
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

	pairMap, err := loadPairMapFromFile(ce.exchanges)
	if err != nil {
		log.Fatal(err)
	}
	ce.positionCrossPairMap = pairMap

	d := time.Duration(time.Millisecond * time.Duration(ce.init.DelayMilliSecond))

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C

		plusMilliSecond := ce.stratStrategy()
		//fmt.Println("d1.5", plusMilliSecond)
		t.Reset(time.Millisecond * time.Duration(ce.init.DelayMilliSecond+plusMilliSecond))
		//fmt.Println("d2", time.Millisecond*time.Duration(ce.init.DelayMilliSecond+plusMilliSecond))
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
	maxHoldVolume := ce.init.MaxHoldVolume
	//計算最大持有量
	maxHoldVolume = math.Min(maxHoldVolume, topCrossPair.GetMaxHoldVolume()*ce.init.MaxHoldeExchangePercent)
	maxHoldVolume = maxHoldVolume * (1.0 - ce.init.MaxHoldBuffer) // 用1%做緩衝

	// 檢查部位是否可以平倉
	hasClose, mp := positionCloseCheck(ce.positionCrossPairMap, crossPairMap, futures, ce.init)
	if hasClose {
		ce.positionCrossPairMap = mp
		savePairMapToFile(ce.positionCrossPairMap)

		//更新交易所內幣別持有量
		ce.updateExchangeWallet(topCrossPair.askExchange.GetName())
		ce.updateExchangeWallet(topCrossPair.bidExchange.GetName())
	}

	currentUSDVolume := getCurrentUSDVolume(ce.positionCrossPairMap)
	execMinTotalValue := topCrossPair.GetMinTotalVolume()
	orderTotalValue := math.Min(maxHoldVolume-currentUSDVolume, execMinTotalValue)

	// 交易判斷
	// 有利可圖
	if !canOrder(maxProfit, orderTotalValue, currentUSDVolume, ce.init) {
		// 無利可圖，重設偵測
		return 0
	}
	log.Println(orderTotalValue)

	m_expectedTotalValue = execMinTotalValue - orderTotalValue

	// 進行交易
	orderCrossPair(topCrossPair, futures, orderTotalValue, ce.init)
	// 調整成交量，改用下單的量，後續平倉成交量才會正確。
	topCrossPair.orderVolume = orderTotalValue
	ce.positionCrossPairMap[topCrossPair.GetName()] = append(ce.positionCrossPairMap[topCrossPair.GetName()], topCrossPair)
	savePairMapToFile(ce.positionCrossPairMap)

	plusMilliSecond := ce.init.DelayMilliSecond
	return plusMilliSecond
}

func (ce *CrossExchange) updateExchangeWallet(exchangeName string) {
	for i, _ := range ce.exchanges {
		if ce.exchanges[i].GetName() == exchangeName {
			ce.exchanges[i].GetWallet()
		}
	}
}

func positionCloseCheck(crossPairsTable map[string][]CrossPair, matchMap CrossPairStringMap, futures exc.Futures, init CrossExchangeInit) (bool, map[string][]CrossPair) {

	//hasPosition
	if len(crossPairsTable) <= 0 {
		return false, crossPairsTable
	}

	hasClose := false

	for key, arrayPairs := range crossPairsTable {

		isClose, arrayPairs := isCloseCheck(key, arrayPairs, matchMap, futures, init)
		if !isClose {
			continue
		}

		hasClose = true

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

	return hasClose, crossPairsTable
}

func getMatchCrossPair(positionPairName string, crossPairArray []CrossPair, matchMap CrossPairStringMap) CrossPair {
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
		log.Println("getMatchCrossPair len(crossPairArray) <= 0 ")
		return CrossPair{}
	}

	matchCrossPair, ok := matchMap[matchName]
	if !ok {
		log.Fatalf("matchMap not found. matchName:%s", matchName)
		return CrossPair{}
	}
	return matchCrossPair
}

func getTotalVolume(crossPairArray []CrossPair, matchCrossPair CrossPair, init CrossExchangeInit) (float64, float64, float64, []CrossPair, string) {
	resSumString := ""
	matchProfitNumber := matchCrossPair.GetProfitNumber()

	matchVolume := matchCrossPair.GetMinTotalVolume()
	// 確認利潤並統計總量
	askTotalVolume := 0.0
	bidTotalVolume := 0.0
	totalMatchOrderUSDVolume := 0.0
	for index, pcp := range crossPairArray {
		positionCrossPair := pcp

		//反向的利潤
		matchProfit := matchProfitNumber / positionCrossPair.GetAskPriceWithFee()

		//找出反向配對，確定利潤
		positionProfit := positionCrossPair.GetProfit()

		log.Println(fmt.Sprintf("position profit:%f; orderVolume:%f, matchProfit:%f", positionProfit, positionCrossPair.orderVolume, matchProfit))
		sum := positionProfit + matchProfit

		hasProfit := (matchProfit > init.MinSellProfit && sum > init.MinSumProfit && matchVolume > init.MinVolume)

		percent := positionCrossPair.GetMinPricePercent(matchCrossPair)
		hasClose := percent < init.StopLosePercent

		if hasProfit || hasClose {
			sumString := fmt.Sprintf("sum:%f, hasProfit:%v, StopPercent:%v, hasClose:%v", sum, hasProfit, percent, hasClose)
			log.Println(sumString)
			resSumString = resSumString + sumString

			thisMatchOrderVolume := math.Min(matchVolume, positionCrossPair.orderVolume)

			askVolume := positionCrossPair.GetAskVolumeByTotal(thisMatchOrderVolume)
			bidVolume := positionCrossPair.GetBidVolumeByTotal(thisMatchOrderVolume)

			askTotalVolume += askVolume
			bidTotalVolume += bidVolume

			totalMatchOrderUSDVolume += thisMatchOrderVolume

			crossPairArray[index].orderVolume -= thisMatchOrderVolume
		}
	}
	return askTotalVolume, bidTotalVolume, totalMatchOrderUSDVolume, crossPairArray, resSumString
}

func isCloseCheck(positionPairName string, crossPairArray []CrossPair, matchMap CrossPairStringMap, futures exc.Futures, init CrossExchangeInit) (bool, []CrossPair) {

	matchCrossPair := getMatchCrossPair(positionPairName, crossPairArray, matchMap)

	askTotalVolume, bidTotalVolume, totalMatchOrderUSDVolume, crossPairArray, sumString := getTotalVolume(crossPairArray, matchCrossPair, init)

	if totalMatchOrderUSDVolume <= 0 {
		return false, crossPairArray
	}

	// 表示有交易
	matchAskExchange, askPair := matchCrossPair.GetAskInfo()
	matchBidExchange, bidPair := matchCrossPair.GetBidInfo()

	isClose := true
	askChannel := executeOrder(matchAskExchange, futures, askPair.Price, exc.Ask, bidTotalVolume, init.OverPrice, isClose)
	bidChannel := executeOrder(matchBidExchange, futures, bidPair.Price, exc.Bid, askTotalVolume, init.OverPrice, isClose)
	//等上面兩個交易都完成，再繼續
	<-askChannel
	<-bidChannel

	content := fmt.Sprintf(
		"positionPairName:%s,\r\n matchPair:%s\r\n orderVolume:%f \r\n sumString:%s",
		positionPairName,
		matchCrossPair.GetProfitString(),
		totalMatchOrderUSDVolume,
		sumString)
	message_tool.SendBroadcastArcherGroup(content)

	return true, crossPairArray
}

func getCrossPairMap(exchanges []exc.Exchange, futures exc.Futures) CrossPairStringMap {

	exchangePairs := []ExchangePair{}
	for _, exchg := range exchanges {
		askPair, bidPair := exchg.GetFuturesAskBidPair(futures)
		exchangePairs = append(exchangePairs, *NewExchangePair(exchg, askPair, bidPair))
	}

	crossPairMap := CrossPairStringMap{}
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

func getMaxProfitCrossPair(mcp CrossPairStringMap) CrossPair {
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

func orderCrossPair(topCrossPair CrossPair, futures exc.Futures, orderTotalValue float64, init CrossExchangeInit) {
	askExchange, askPair := topCrossPair.GetAskInfo()
	bidExchange, bidPair := topCrossPair.GetBidInfo()

	askVolume := askExchange.GetVolumeByTotal(orderTotalValue, askPair.Price)
	bidVolume := bidExchange.GetVolumeByTotal(orderTotalValue, bidPair.Price)

	isClose := false
	askChannel := executeOrder(askExchange, futures, askPair.Price, exc.Ask, askVolume, init.OverPrice, isClose)
	bidChannel := executeOrder(bidExchange, futures, bidPair.Price, exc.Bid, bidVolume, init.OverPrice, isClose)
	//等上面兩個交易都完成，再繼續
	<-askChannel
	<-bidChannel

	content := fmt.Sprintf("%s, %s\r\n orderTotalValue:%g \r\n maxProfit:%g \r\n m_expectedTotalValue:%g",
		fmt.Sprintf("resAsk:%f, orderTotal:%f, AskCoin:%s", askPair.Price, askPair.Total(), askExchange.GetName()),
		fmt.Sprintf("resBid:%f, orderTotal:%f, bidCoin:%s", bidPair.Price, bidPair.Total(), bidExchange.GetName()),
		orderTotalValue,
		topCrossPair.GetProfit(),
		m_expectedTotalValue)
	log.Println(content)

	message_tool.SendBroadcastArcherGroup(content)
}

var (
	m_expectedTotalValue float64 = 0
)

const (
	PairMapFileName = "PairMap.json"
)

func getCurrentUSDVolume(positionCrossPairMap map[string][]CrossPair) float64 {
	sum := 0.0
	for _, array := range positionCrossPairMap {
		for _, v := range array {
			sum = sum + v.orderVolume
		}
	}
	return sum
}

func canOrder(profit, orderTotalValue, currentUSDVolume float64, init CrossExchangeInit) bool {
	// 有利可圖
	if profit < init.MinCreateProfit { // 沒足夠利潤，直接下一圈
		log.Println(fmt.Sprintf("No enough profit. profit:%f", profit))
		return false
	} else if orderTotalValue < init.MinVolume {
		log.Println(fmt.Sprintf("orderTotalValue:%f < %f", orderTotalValue, init.MinVolume))
		return false
	} else if currentUSDVolume >= init.MaxHoldVolume {
		log.Println(fmt.Sprintf("currentUSDVolume:%g >= init.MaxHoldVolume:%g", currentUSDVolume, init.MaxHoldVolume))
		return false
	}
	return true
}

func executeOrder(exchange exc.Exchange, futures exc.Futures, price float64, pType exc.PriceType, volume float64, overPrice float64, isClose bool) chan int {
	resultChannel := make(chan int)
	side := "sell"
	adjPrice := price * (1.0 - overPrice)
	if pType == exc.Ask {
		side = "buy"
		adjPrice = price * (1.0 + overPrice)
	}

	myOrder := exc.FuturesOrder{
		CommodityOrder: exc.CommodityOrder{
			Side:      side,
			Price:     adjPrice,
			Size:      volume,
			OrderType: exc.LIMIT,
		},
		Futures: futures,
		IsClose: isClose,
	}
	go func() {
		exc.PostFuturesOrderRefry(exchange, myOrder)
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
