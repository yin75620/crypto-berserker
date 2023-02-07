package CrossExchange

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/yin75620/crypto-berserker/ksql"

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
	sql                  *ksql.Ksql
	mutex                sync.Mutex
}

var (
	mPreWallets []exc.Wallet
)

func NewCrossExchange(exchanges []exc.Exchange) *CrossExchange {
	ce := CrossExchange{}
	ce.exchanges = exchanges
	ce.positionCrossPairMap = map[string][]CrossPair{}
	NewCrossExchangeInit()
	ce.init = *NewCrossExchangeInit()
	ce.sql = ksql.NewKsql()

	mPreWallets = make([]exc.Wallet, len(ce.exchanges))
	return &ce
}

//func (ce *CrossExchange) SetInit(init CrossExchangeInit) {
//	ce.init = init
//}

func (ce *CrossExchange) SetInitByIni(filename string) {
	ce.init.IniSetting(filename)
}

func (ce *CrossExchange) SetInitExtraProfitByIni(filename string, futureNames []string) {
	ce.init.InitExtraProfit(filename, futureNames)
}

func (ce *CrossExchange) SetFuturesArray(futuresArray []exc.Futures) {
	ce.futuresArray = futuresArray
}

func (ce *CrossExchange) Start() {
	slog := simpleLog.StartLog()
	defer slog.Close()

	message_tool.StartTelegram()
	if ce.init.IsEnableDBLog {
		err := ce.sql.Start()
		if err != nil {
			log.Println("ERROR: Database not connect.", err)

		}
		defer ce.sql.End()
	}

	//test api
	infoAll := "Start"
	for i, exchg := range ce.exchanges {
		exchg.GetAccountInfo()
		wallet := exchg.GetWallet()
		infoAll = fmt.Sprintf("%s \r\n %s: %v", infoAll, exchg.GetName(), wallet)
		mPreWallets[i] = wallet
	}
	log.Println(string(infoAll))
	message_tool.SendBroadcastArcherGroup(infoAll)

	pairMap, err := loadPairMapFromFile(ce.exchanges)
	if err != nil {
		log.Fatal(err)
	}
	ce.positionCrossPairMap = pairMap

	liveTime := time.Second * time.Duration(ce.init.ShowLiveSecond)
	liveTimer := time.NewTimer(liveTime)

	ce.stratStrategy()

	for {
		select {
		case <-liveTimer.C:
			log.Println("live")
			liveTimer.Reset(liveTime)
		}

	}

}

func (ce *CrossExchange) stratStrategy() int64 {

	var totalWaitTime int64
	for _, futures := range ce.futuresArray {
		mf := futures

		go func() {

			d := time.Duration(time.Millisecond * time.Duration(ce.init.DelayMilliSecond))
			t := time.NewTimer(d)
			defer t.Stop()

			for {
				select {
				case <-t.C:

					plusMilliSecond := ce.stratFuturesStrategy(mf)
					t.Reset(time.Millisecond * time.Duration(ce.init.DelayMilliSecond+plusMilliSecond))
				}
			}

		}()

	}
	return totalWaitTime
}

func (ce *CrossExchange) stratFuturesStrategy(futures exc.Futures) int64 {

	crossPairMap := getCrossPairMap(ce.exchanges, futures)

	// 沒有足夠的交易所可以套利
	if len(crossPairMap) == 0 {
		log.Fatal("has no match pair exchange")
	}

	topCrossPair := ce.getMaxProfitCrossPair(crossPairMap)

	if isDisconnect(topCrossPair) {
		var disconnectMiliSeccond int64 = 3000
		log.Println(fmt.Sprintf("error: disconnect, wait... %v miliSecond.", disconnectMiliSeccond))
		return disconnectMiliSeccond
	}

	// 檢查部位是否可以平倉
	hasClose, mp := ce.positionCloseCheck(ce.positionCrossPairMap, crossPairMap, futures, ce.init)
	if hasClose {
		ce.positionCrossPairMap = mp
		savePairMapToFile(ce.positionCrossPairMap)

		//完成N次交易報告資產變化值
		const EveryCountCheckWallet = 1
		var mFinishCount = 1
		if mFinishCount%EveryCountCheckWallet == 0 {
			sendInfo := ""
			ce.mutex.Lock()
			for i, exchange := range ce.exchanges {
				wallet := exchange.GetWallet()
				array := mPreWallets[i].GetAllBalanceProfit(wallet)
				mPreWallets[i] = wallet
				sendInfo = fmt.Sprintf("%s%s", sendInfo, fmt.Sprintf("%s: %v, ", exchange.GetName(), array))
			}
			ce.mutex.Unlock()
			log.Println(sendInfo)
			message_tool.SendBroadcastArcherGroup(sendInfo)
		}
		return ce.init.DealtDelayMilliSecond //成交之後要等待多久
	}

	topCrossPair = ce.getCanOrderMaxProfitCrossPair(crossPairMap, futures)
	if topCrossPair.IsDefault() {
		// 無利可圖，重設偵測
		return 0
	}

	// 交易判斷
	// 有利可圖
	isCaneOrder, orderTotalValue, execMinTotalValue := ce.coinPairCanOrder(topCrossPair, futures)
	if !isCaneOrder {
		// 無利可圖，重設偵測
		return 0
	}
	log.Println(orderTotalValue)

	m_expectedTotalValue = execMinTotalValue - orderTotalValue

	ce.mutex.Lock()
	// 進行交易
	orderCrossPair(topCrossPair, futures, orderTotalValue, ce.init)
	//更新交易所內幣別持有量
	ce.updateExchangeWallet(topCrossPair.askExchange.GetName())
	ce.updateExchangeWallet(topCrossPair.bidExchange.GetName())
	ce.mutex.Unlock()

	// 調整成交量，改用下單的量，後續平倉成交量才會正確。
	topCrossPair.orderVolume = orderTotalValue
	ce.positionCrossPairMap[topCrossPair.GetName()] = append(ce.positionCrossPairMap[topCrossPair.GetName()], topCrossPair)
	savePairMapToFile(ce.positionCrossPairMap)

	return ce.init.DealtDelayMilliSecond //成交之後要等待多久
}

func isDisconnect(cp CrossPair) bool {
	return cp.askPricePair.Price == 0 || cp.askPricePair.Volume == 0 || cp.bidPricePair.Price == 0 || cp.bidPricePair.Volume == 0
}

func (ce *CrossExchange) updateExchangeWallet(exchangeName string) {
	for i, _ := range ce.exchanges {
		if ce.exchanges[i].GetName() == exchangeName {
			ce.exchanges[i].GetWallet()
			ce.exchanges[i].GetAccountInfo() // 更新未平倉損億
		}
	}
}

func (ce *CrossExchange) positionCloseCheck(crossPairsTable map[string][]CrossPair, matchMap CrossPairStringMap, futures exc.Futures, init CrossExchangeInit) (bool, map[string][]CrossPair) {

	//hasPosition
	if len(crossPairsTable) <= 0 {
		return false, crossPairsTable
	}

	sameFutures := make(map[string][]CrossPair)
	diffFutures := make(map[string][]CrossPair)
	// 拆成想要這次的貨幣與非這次的貨幣
	for key, arrayPairs := range crossPairsTable {
		for _, v := range arrayPairs {
			if v.Symbol == futures.GetMarketName() {
				sameFutures[key] = append(sameFutures[key], v)
			} else {
				diffFutures[key] = append(diffFutures[key], v)
			}
		}
	}

	hasClose := false

	for key, arrayPairs := range sameFutures {

		isClose, arrayPairs := ce.isCloseCheck(key, arrayPairs, matchMap, futures, init)
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
			delete(sameFutures, key)
		} else {
			sameFutures[key] = arrayPairs
		}
	}

	//再組合
	for k, v := range sameFutures {
		diffFutures[k] = append(diffFutures[k], v...)
	}
	crossPairsTable = diffFutures

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

func getTotalVolume(crossPairArray []CrossPair, matchCrossPair CrossPair, init CrossExchangeInit, futures exc.Futures) (float64, float64, float64, []CrossPair, string) {
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

		//log.Println(fmt.Sprintf("position profit:%f; orderVolume:%f, matchProfit:%f", positionProfit, positionCrossPair.orderVolume, matchProfit))
		sum := positionProfit + matchProfit

		hasProfit := (matchProfit > getMinSellProfit(init, futures) && sum > getMinSumProfit(init, futures) && matchVolume > init.MinVolume)

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

func (ce *CrossExchange) isCloseCheck(positionPairName string, crossPairArray []CrossPair, matchMap CrossPairStringMap, futures exc.Futures, init CrossExchangeInit) (bool, []CrossPair) {

	matchCrossPair := getMatchCrossPair(positionPairName, crossPairArray, matchMap)

	askTotalVolume, bidTotalVolume, totalMatchOrderUSDVolume, crossPairArray, sumString := getTotalVolume(crossPairArray, matchCrossPair, init, futures)

	if totalMatchOrderUSDVolume <= 0 {
		return false, crossPairArray
	}

	speedTestStart := time.Now()
	ce.mutex.Lock()
	// 表示有交易
	matchAskExchange, askPair := matchCrossPair.GetAskInfo()
	matchBidExchange, bidPair := matchCrossPair.GetBidInfo()

	isClose := true
	askChannel := executeOrder(matchAskExchange, futures, askPair.Price, exc.Ask, bidTotalVolume, init.OverPrice, isClose)
	bidChannel := executeOrder(matchBidExchange, futures, bidPair.Price, exc.Bid, askTotalVolume, init.OverPrice, isClose)
	//等上面兩個交易都完成，再繼續
	<-askChannel
	<-bidChannel

	elapsed := time.Since(speedTestStart)

	//更新交易所內幣別持有量
	ce.updateExchangeWallet(matchAskExchange.GetName())
	ce.updateExchangeWallet(matchBidExchange.GetName())
	ce.mutex.Unlock()

	content := fmt.Sprintf(
		"%s positionPairName: %s,\r\nmatchPair: %s\r\norderVolume: %f \r\nsumString: %s\r\ntime: %s\r\nelapsed: %v",
		futures.GetMarketName(),
		positionPairName,
		matchCrossPair.GetProfitString(),
		totalMatchOrderUSDVolume,
		sumString,
		time.Now().UTC(),
		elapsed)

	go func() {
		//wait then send to telegram
		time.Sleep(time.Millisecond * time.Duration(ce.init.DealtDelayMilliSecond))
		//看一下交易紀錄
		lastTradeInfo := getLastTradeInfo(matchAskExchange, matchBidExchange, futures)
		sendContent := fmt.Sprintf("%s \r\nlastTradeInfo: %s", content, lastTradeInfo)
		message_tool.SendBroadcastArcherGroup(sendContent)
	}()

	log.Println(content)

	return true, crossPairArray
}

func getLastTradeInfo(askExchange exc.Exchange, bidExchange exc.Exchange, futures exc.Futures) string {
	symbol := futures.GetLinkMarketNameUpper()
	askTrades := askExchange.GetTightUserTrades(symbol)
	bidTrades := bidExchange.GetTightUserTrades(symbol)

	askItem := getFirstItme(askTrades)
	bidItem := getFirstItme(bidTrades)

	return fmt.Sprintf("%s:%s %f v:%f t:%v \r\n%s:%s %f v:%f t:%v",
		askExchange.GetName(), askItem.Side, askItem.Price, askItem.Qty, askItem.Time,
		bidExchange.GetName(), bidItem.Side, bidItem.Price, bidItem.Qty, bidItem.Time)
}

func getFirstItme(uts map[exc.UserTradeKey]exc.UserTrade) exc.UserTrade {

	for _, v := range uts {
		return v
	}
	return exc.UserTrade{}
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
			cp := NewCrossPair(ePair.exchange, matchPair.exchange, ePair.ask, matchPair.bid, futures.GetMarketName())
			crossPairMap[cp.GetName()] = *cp
		}
	}
	return crossPairMap
}

func (ce *CrossExchange) getMaxProfitCrossPair(mcp CrossPairStringMap) CrossPair {
	maxProfit := -math.MaxFloat64
	var topCrossPair CrossPair
	for _, crossPair := range mcp {
		if crossPair.GetProfit() > maxProfit {
			maxProfit = crossPair.GetProfit()
			topCrossPair = crossPair
		}

		if ce.init.IsEnableDBLog {
			s := crossPair.GetDbInserString()
			go func() {
				err := ce.sql.Insert(s)
				if err != nil {
					log.Println("ce.sql.Insert", err)
				}
			}()
		} else {
			if ce.init.IsShowProfitLog {
				s := crossPair.GetProfitString()
				log.Println(s)
			}

		}

		//log.Println(crossPair.GetProfitString())
		//matchPair := crossPairMap[crossPair.GetMatchName()]
		//totalprofit := matchPair.GetProfit() + crossPair.GetProfit()
	}
	return topCrossPair
}

func (ce *CrossExchange) getCanOrderMaxProfitCrossPair(mcp CrossPairStringMap, futures exc.Futures) CrossPair {
	maxProfit := -math.MaxFloat64
	var topCrossPair CrossPair
	for _, crossPair := range mcp {
		canOrder, _, _ := ce.coinPairCanOrder(crossPair, futures)
		if !canOrder {
			continue
		}

		if crossPair.GetProfit() > maxProfit {
			maxProfit = crossPair.GetProfit()
			topCrossPair = crossPair
		}
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

	content := fmt.Sprintf("%s, %s, %s\r\norderTotalValue:%g \r\nmaxProfit:%g \r\nexpectedTotalValue:%g",
		futures.GetMarketName(),
		fmt.Sprintf("resAsk:%f, orderTotal:%f, AskCoin:%s", askPair.Price, askPair.Total(), askExchange.GetName()),
		fmt.Sprintf("resBid:%f, orderTotal:%f, bidCoin:%s", bidPair.Price, bidPair.Total(), bidExchange.GetName()),
		orderTotalValue,
		topCrossPair.GetProfit(),
		m_expectedTotalValue)
	log.Println(content)

	go func() {
		//wait 2 second then send to telegram
		time.Sleep(time.Millisecond * time.Duration(2000))
		//看一下交易紀錄
		lastTradeInfo := getLastTradeInfo(askExchange, bidExchange, futures)
		sendContent := fmt.Sprintf("%s \r\nlastTradeInfo: %s", content, lastTradeInfo)
		message_tool.SendBroadcastArcherGroup(sendContent)
	}()
}

var (
	m_expectedTotalValue float64 = 0
)

const (
	PairMapFileName = "PairMap.json"
)

func getCurrentUSDVolume(name string, positionCrossPairMap map[string][]CrossPair) float64 {
	sum := 0.0
	for keyName, array := range positionCrossPairMap {
		if keyName != name {
			continue
		}
		for _, v := range array {
			sum = sum + v.orderVolume
		}
	}
	return sum
}

func canOrder(profit, orderTotalValue, currentUSDVolume float64, init CrossExchangeInit, futures exc.Futures) bool {
	// 有利可圖
	if profit < getMinCreateProfit(init, futures) { // 沒足夠利潤，直接下一圈
		//log.Println(fmt.Sprintf("No enough profit. profit:%f", profit))
		return false
	} else {
		log.Println(fmt.Sprintf("Futures:%s,  profit:%f", futures.GetMarketName(), profit))
		if orderTotalValue < init.MinVolume {
			log.Println(fmt.Sprintf("orderTotalValue:%f < init.MinVolume:%f", orderTotalValue, init.MinVolume))
			return false
		} else if currentUSDVolume >= init.MaxHoldVolume {
			log.Println(fmt.Sprintf("currentUSDVolume:%g >= init.MaxHoldVolume:%g", currentUSDVolume, init.MaxHoldVolume))
			return false
		}
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

func (ce *CrossExchange) getCalculateMaxHoldVolume(crossPair CrossPair) float64 {
	maxHoldVolume := ce.init.MaxHoldVolume
	//計算最大持有量
	ce.mutex.Lock()
	maxHoldVolume = math.Min(maxHoldVolume, crossPair.GetMaxHoldVolume()*ce.init.MaxHoldeExchangePercent)
	ce.mutex.Unlock()
	maxHoldVolume = maxHoldVolume * (1.0 - ce.init.MaxHoldBuffer) // 用1%做緩衝
	return maxHoldVolume
}

func (ce *CrossExchange) coinPairCanOrder(crossPair CrossPair, futures exc.Futures) (bool, float64, float64) {
	maxProfit := crossPair.GetProfit()
	maxHoldVolume := ce.getCalculateMaxHoldVolume(crossPair)
	currentUSDVolume := getCurrentUSDVolume(crossPair.GetName(), ce.positionCrossPairMap)
	execMinTotalValue := crossPair.GetMinTotalVolume()
	orderTotalValue := math.Min(maxHoldVolume-currentUSDVolume, execMinTotalValue)

	return canOrder(maxProfit, orderTotalValue, currentUSDVolume, ce.init, futures), orderTotalValue, execMinTotalValue
}

func getMinSellProfit(cei CrossExchangeInit, futures exc.Futures) float64 {
	if v, ok := cei.ExtraFutures[futures.GetIniNameUpper()]; ok {
		return v.MinSellProfit
	}
	return cei.MinSellProfit
}

func getMinSumProfit(cei CrossExchangeInit, futures exc.Futures) float64 {
	if v, ok := cei.ExtraFutures[futures.GetIniNameUpper()]; ok {
		return v.MinSumProfit
	}
	return cei.MinSumProfit
}

func getMinCreateProfit(cei CrossExchangeInit, futures exc.Futures) float64 {
	if v, ok := cei.ExtraFutures[futures.GetIniNameUpper()]; ok {
		return v.MinCreateProfit
	}
	return cei.MinCreateProfit
}
