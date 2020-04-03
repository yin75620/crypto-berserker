package Triangular

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
	"github.com/yin75620/crypto-berserker/jmath"
	"github.com/yin75620/crypto-berserker/message_tool"
)

const (
	USD  = "USD"
	USDT = "USDT"
	USDC = "USDC"
	BTC  = "BTC"
	FTT  = "FTT"
	PAX  = "PAX"
	BTMX = "BTMX"
	ETH  = "ETH"
	LTC  = "LTC"
	XRP  = "XRP"
)

var marketMap map[string]int = map[string]int{"BTC/USD": 4}

const (
	R_PROFIT      = 0
	R_PLUS_SECOND = 1
	R_TOTAL_VALUE = 2
)

type TriangularInit struct {
	RangePremium    float64
	LeastTotalValue float64
	DelayTime       int
	// 利潤, 加速, 價值(美元計價)
	RANK_S []float64 //{0.006, -3.0, 1000}
	RANK_N []float64 //{0.003, -2.0, 300}
}

type CoinStrip struct {
	Coins []string
}

type CoinBunch struct {
	CoinStrips []CoinStrip
}

type Triangular struct {
	exchangeClient exc.Exchange
	CoinBunch      []CoinBunch
	Init           TriangularInit
}

func NewTriangular(exchange exc.Exchange) *Triangular {
	var t = &Triangular{
		Init: TriangularInit{
			RangePremium:    0.1, //10%
			LeastTotalValue: 10,  //10 quoteCoin
			DelayTime:       5,   //5second
			// 利潤, 加速, 價值(美元計價)
			RANK_S: []float64{0.006, -3.0, 1000},
			RANK_N: []float64{0.001, -2.0, 300}}}
	t.exchangeClient = exchange
	return t
}

func (tri *Triangular) SetCoinBunchs(coinBunchs []CoinBunch) {
	for _, value := range coinBunchs {
		tri.CoinBunch = append(tri.CoinBunch, value)
	}
}

func (tri *Triangular) SetInit(init TriangularInit) {
	tri.Init = init
}

func (tri *Triangular) Start() {
	message_tool.StartTelegram()
	var logFile *os.File = StartLog()
	defer logFile.Close()
	tri.stratStrategy()

	infoStr := string(tri.exchangeClient.GetAccountInfo())
	message_tool.SendTelegram(infoStr)

	mPreWallet = tri.exchangeClient.GetWallet()

	var delay_time int = tri.Init.DelayTime
	d := time.Duration(time.Second * time.Duration(delay_time))

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C
		plusSecond := tri.stratStrategy()
		t.Reset(time.Second * time.Duration(delay_time+plusSecond))
	}

	///開啟伺服器讓程式留著
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	///
}

func StartLog() *os.File {
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))

	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	return logFile
}

type Quote struct {
	goalCoin    string
	currentCoin string
	askPair     exc.PricePair
	bidPair     exc.PricePair
}

func (q *Quote) GetPair(pType exc.PriceType) exc.PricePair {
	switch pType {
	case exc.Ask:
		return q.askPair
	case exc.Bid:
		return q.bidPair
	}
	log.Fatal("Error: not specific pType")
	return exc.PricePair{}
}

func (q *Quote) GetCoinPair() exc.CoinPair {
	res := exc.CoinPair{BaseCoin: q.goalCoin, QuotedCoin: q.currentCoin}
	return res
}

func (tri *Triangular) NewQuote(goalCoin, currentCoin string) Quote {
	var askPair exc.PricePair
	var bidPair exc.PricePair
	askPair, bidPair = tri.exchangeClient.GetAskBidPair(exc.CoinPair{BaseCoin: goalCoin, QuotedCoin: currentCoin}, 1)
	var quote Quote
	quote.askPair = askPair
	quote.bidPair = bidPair
	quote.goalCoin = goalCoin
	quote.currentCoin = currentCoin
	return quote
}

type DealFlow struct {
	quotes   []Quote
	takerFee float64
}

type QuoteChannelRes struct {
	index int
	quote Quote
}

func (tri *Triangular) NewDealFlow(goalCoin string, stepCoins []string) DealFlow {
	dealFlow := DealFlow{}
	dealFlow.takerFee = tri.exchangeClient.GetFee().Taker

	tempGoalCoin := goalCoin
	stepLen := len(stepCoins)
	dealFlow.quotes = make([]Quote, stepLen)

	finishChannel := make(chan QuoteChannelRes, stepLen)
	for i := 0; i < stepLen; i++ {
		stepCoin := stepCoins[i]
		go func(tempGoalCoin string, stepCoin string, index int) {
			quote := tri.NewQuote(tempGoalCoin, stepCoin)
			//dealFlow.quotes[i] = quote
			//tempGoalCoin = stepCoin
			finishChannel <- QuoteChannelRes{index: index, quote: quote}
		}(tempGoalCoin, stepCoin, i)
		tempGoalCoin = stepCoin
	}

	// 等待全部完成再回傳
	for i := 0; i < stepLen; i++ {
		res := <-finishChannel
		dealFlow.quotes[res.index] = res.quote
	}

	return dealFlow
}

func (df *DealFlow) getName() string {
	marketName := ""
	for _, value := range df.quotes {
		co := value.GetCoinPair()
		marketName = fmt.Sprintf("%s%s", marketName, co.GetMarketName())
	}
	return marketName
}

func (df *DealFlow) getFinalPair(pType exc.PriceType) exc.PricePair {
	return df.getFinalPairWithFee(pType, true)
}

func (df *DealFlow) getFinalPairWithFee(pType exc.PriceType, hasFee bool) exc.PricePair {
	var finalPrice float64 = 1
	var finalFeePrice float64 = 1
	var resTotalValue float64 = math.MaxFloat64 //總價值

	for _, quote := range df.quotes {
		pair := quote.GetPair(pType)
		finalPrice = finalPrice * pair.Price
		finalFeePrice = finalFeePrice * pair.Price
		if hasFee {
			cacularSymbol := 1.0
			if pType == exc.Bid {
				cacularSymbol = -1
			}
			finalFeePrice = finalFeePrice * (1.0 + df.takerFee*cacularSymbol)
		}

		// 上次總價值 與這次量相比 (這樣單位才一致)
		minVolume := math.Min(resTotalValue, pair.Volume)

		resTotalValue = pair.Price * minVolume
	}

	var finalAskPair exc.PricePair = exc.PricePair{}
	finalAskPair.Price = finalFeePrice
	finalAskPair.Volume = resTotalValue / finalPrice
	return finalAskPair
}
func getLowestFlow(dealFlows []DealFlow, pType exc.PriceType) DealFlow {
	lowest := math.MaxFloat64
	resDealFlow := DealFlow{}
	for _, value := range dealFlows {
		pair := value.getFinalPair(pType)
		log.Println(fmt.Sprintf("getLowestFlow:%f, Coin:%s", pair.Price, value.getName()))
		if lowest > pair.Price {
			lowest = pair.Price
			resDealFlow = value
		}
	}
	return resDealFlow
}

func getHighestFlow(dealFlows []DealFlow, pType exc.PriceType) DealFlow {
	highest := 0.0
	resDealFlow := DealFlow{}
	for _, value := range dealFlows {
		pair := value.getFinalPair(pType)
		log.Println(fmt.Sprintf("getHighestFlow:%f, Coin:%s", pair.Price, value.getName()))
		if highest < pair.Price {
			highest = pair.Price
			resDealFlow = value
		}
	}
	return resDealFlow
}

var (
	m_expectedTotalValue   float64 = 0
	m_expectedLowestProfit float64 = 0
)

var m_isFullPower = false

var mFailCount = 0
var MAX_FAIL_COUNT = 3

var mFinishCount = 0

const EveryCountCheckWallet = 5
const DefaultDelaySecond = 100

var mPreWallet exc.Wallet

func (tri *Triangular) stratStrategy() int {
	res := DefaultDelaySecond
	for _, coinBunch := range tri.CoinBunch {
		dealFlows := []DealFlow{}
		for _, coinStrip := range coinBunch.CoinStrips {
			fuDealFlow := tri.NewDealFlow(coinStrip.Coins[0], coinStrip.Coins[1:])
			dealFlows = append(dealFlows, fuDealFlow)
		}
		sec := tri.stratDealFlow(dealFlows)
		if sec < res {
			res = sec
		}
	}
	return res
}

func (tri *Triangular) stratDealFlow(dealFlows []DealFlow) int {

	defer func() {
		err := recover()
		if err == nil {
			return
		}
		log.Println(err) // 這已經是頂層的 UI 介面了，想以自己的方式呈現錯誤
		mFailCount = mFailCount + 1
		if mFailCount < MAX_FAIL_COUNT {
			//再重來一次
			tri.stratDealFlow(dealFlows)
		}
		log.Println(mFailCount)
		// 失敗次數太多，直接結束
	}()

	/*dealFlows := []DealFlow{}

	for _, coinStrip := range tri.CoinStrip {
		fuDealFlow := tri.NewDealFlow(coinStrip.Coins[0], coinStrip.Coins[1:])
		dealFlows = append(dealFlows, fuDealFlow)
	}*/

	lowestAskFlow := getLowestFlow(dealFlows, exc.Ask)
	highestBidFlow := getHighestFlow(dealFlows, exc.Bid)

	laName := lowestAskFlow.getName()
	hbName := highestBidFlow.getName()

	laPrice := lowestAskFlow.getFinalPair(exc.Ask).Price
	hbPrice := highestBidFlow.getFinalPair(exc.Bid).Price

	laVolume := lowestAskFlow.getFinalPair(exc.Ask).Volume
	hbVolume := highestBidFlow.getFinalPair(exc.Bid).Volume
	// 出現錯誤，放慢速度
	if laPrice <= 0 {
		log.Println("laPrice <= 0")
		return 60
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

	currentOrderTotalValue := tri.Init.RANK_N[R_TOTAL_VALUE]

	// 表示有人來搶單拉!!
	if m_expectedTotalValue != 0 && m_expectedTotalValue > currentOrderTotalValue {
		m_isFullPower = true

	} else if profit > tri.Init.RANK_S[R_PROFIT] {
		// 利潤超高 買起來!!!
		m_isFullPower = true
	} else if m_expectedLowestProfit != 0 && m_expectedLowestProfit < profit {
		// 利潤變少了，全力買起來
		m_isFullPower = true
	}

	if m_isFullPower {
		currentOrderTotalValue = tri.Init.RANK_S[R_TOTAL_VALUE]
	}

	wnatOrderTotalValue := currentOrderTotalValue * (1 - (10 * rand.Float64() / 100.0)) // 隨機 -10%
	wnatOrderTotalValue = math.Floor(wnatOrderTotalValue)

	orderTotalValue := math.Min(wnatOrderTotalValue, minSourceTotalValue)

	m_expectedTotalValue = minSourceTotalValue - orderTotalValue
	m_expectedLowestProfit = profit

	// 有利可圖
	if !tri.canOrder(profit, orderTotalValue) {
		// 無利可圖，重設偵測
		m_isFullPower = false
		m_expectedTotalValue = 0
		m_expectedLowestProfit = 0
		return 0
	}

	laOrderVolume := orderTotalValue / laPrice
	hbOrderVolume := orderTotalValue / hbPrice
	const isOrder = true
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
	}

	content := fmt.Sprintf("%s, %s\r\n %s,\r\n orderTotalValue:%g \r\n profit:%g \r\n m_expectedTotalValue:%g",
		tri.exchangeClient.GetName(),
		fmt.Sprintf("resAsk:%f, orderVolume:%f, AskCoin:%s", laPrice, laOrderVolume, laName),
		fmt.Sprintf("resBid:%f, orderVolume:%f, bidCoin:%s", hbPrice, hbOrderVolume, hbName),
		orderTotalValue,
		profit,
		m_expectedTotalValue)
	message_tool.SendTelegram(content)

	mFinishCount = mFinishCount + 1
	if mFinishCount%EveryCountCheckWallet == 0 {
		wallet := tri.exchangeClient.GetWallet()
		//完成N次交易報告資產變化值
		array := mPreWallet.GetAllBalanceProfit(wallet)
		sendInfo := fmt.Sprintf("earn:%v", array)
		log.Println(sendInfo)
		message_tool.SendTelegram(sendInfo)

		if mPreWallet.IsAllBalanceReduce(wallet) {
			sendContent := fmt.Sprintf("sustained losses.after %d times", EveryCountCheckWallet)
			message_tool.SendTelegram(sendContent)
			log.Fatal(sendContent)
		}
		mPreWallet = wallet
	}

	resPlusSecond := tri.Init.RANK_N[R_PLUS_SECOND]
	if m_isFullPower {
		resPlusSecond = tri.Init.RANK_S[R_PLUS_SECOND]
	}

	return int(resPlusSecond)
}

func (tri *Triangular) canOrder(profit, orderTotalValue float64) bool {
	// 有利可圖
	if profit < 0 {
		log.Println("No profit")
		return false
	} else if profit < tri.Init.RANK_N[R_PROFIT] {
		log.Println("No enough profit")
		return false
	} else if orderTotalValue < tri.Init.LeastTotalValue {
		log.Println(fmt.Sprintf("orderTotalValue < %f", tri.Init.LeastTotalValue))
		return false
	}
	return true
}

///

func (tri *Triangular) executeOrder(df DealFlow, pType exc.PriceType, startVolume float64) chan int {
	log.Println(fmt.Sprintf("startVolume:%f", startVolume))

	allFinishChannel := make(chan int)

	side := ""
	orderSymbol := 1.0
	switch pType {
	case exc.Bid:
		side = "sell"
		orderVolume := startVolume
		orderSymbol = -1
		finishChannel := make(chan int, len(df.quotes))
		for _, quote := range df.quotes {
			var merketInfo = tri.exchangeClient.GetMarketInfo(quote.GetCoinPair())
			unit := merketInfo.VolumeIncrement
			orderVolume = jmath.FloatFloorByFloat(orderVolume, unit)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*tri.Init.RangePremium)

			var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: exc.LIMIT,
				CoinPair:  quote.GetCoinPair(),
			}
			go func() {
				tri.PostOrderRefry(myOrder)
				finishChannel <- 0
			}()

			orderVolume = orderVolume * quote.GetPair(pType).Price
		}

		go func() {
			// 等待完成
			for i := 0; i < len(df.quotes); i++ {
				<-finishChannel
			}
			allFinishChannel <- 0
		}()

	case exc.Ask:
		side := "buy"
		orderVolume := startVolume
		orderSymbol = 1.0
		var orders []exc.ExchangeOrder = []exc.ExchangeOrder{}
		finishChannel := make(chan int, len(df.quotes))
		for i := 0; i < len(df.quotes); i++ {
			quote := df.quotes[i]
			var merketInfo = tri.exchangeClient.GetMarketInfo(quote.GetCoinPair())
			unit := merketInfo.VolumeIncrement
			orderVolume = jmath.FloatFloorByFloat(orderVolume, unit)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*tri.Init.RangePremium)

			var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: exc.LIMIT,
				CoinPair:  quote.GetCoinPair(),
			}
			orders = append(orders, myOrder)

			orderVolume = orderVolume * quote.GetPair(pType).Price

		}

		for i := len(orders) - 1; i >= 0; i-- {
			order := orders[i]
			go func() {
				tri.PostOrderRefry(order)
				finishChannel <- 0
			}()
		}

		go func() {
			//等待完成
			for i := len(orders) - 1; i >= 0; i-- {
				<-finishChannel
			}
			allFinishChannel <- 0
		}()

	}

	return allFinishChannel
}

func (tri *Triangular) PostOrderRefry(order exc.ExchangeOrder) {
	exc.PostOrderRefry(tri.exchangeClient, order)
}
