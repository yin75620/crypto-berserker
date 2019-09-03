package Triangular

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	exc "github.com/yin75620/crypto-berserker/exchange"
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

var (
	mRangePremium    float64 = 0.1 //10%
	mLeastTotalValue float64 = 10  //10US
	mDelayTimeSec    int     = 5
)

//當const 用
var (
	// 數量, 利潤, 加速, 價值(美元計價)
	RANK_S = []float64{0.006, -3.0, 1000}
	RANK_N = []float64{0.001, -2.0, 300}
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

type Triangular struct {
	exchangeClient exc.Exchange
	CoinStrip      []CoinStrip
}

func NewTriangular(exchange exc.Exchange) *Triangular {
	var t = &Triangular{}
	t.exchangeClient = exchange
	return t
}

func (tri *Triangular) SetCoinBunch(CoinStrip []CoinStrip) {
	for _, value := range CoinStrip {
		tri.CoinStrip = append(tri.CoinStrip, value)
	}
}

func (tri *Triangular) SetInit(init TriangularInit) {
	mRangePremium = init.RangePremium       //10%
	mLeastTotalValue = init.LeastTotalValue //10US
	mDelayTimeSec = init.DelayTime
	RANK_N = init.RANK_N
	RANK_S = init.RANK_S
}

func (tri *Triangular) Start() {
	message_tool.StartTelegram()
	var logFile *os.File = StartLog()
	defer logFile.Close()
	tri.stratStrategy()

	infoStr := string(tri.exchangeClient.GetAccountInfo())
	message_tool.SendTelegram(infoStr)

	var delay_time int = mDelayTimeSec
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
	logFile, err := os.OpenFile("earn.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	underDot    int
}

func (q *Quote) MarketName() string {
	marketName := fmt.Sprintf("%s/%s", q.goalCoin, q.currentCoin)
	return marketName
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

func (tri *Triangular) NewQuote(goalCoin, currentCoin string) Quote {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)
	var askPair exc.PricePair
	var bidPair exc.PricePair
	//偽裝成 USDT
	if marketName == "USDT/USD" {
		askPair = exc.PricePair{1.001, 999999}
		bidPair = exc.PricePair{0.997, 999999}
	} else {
		askPair, bidPair = tri.exchangeClient.GetAskBidPair(exc.CoinPair{BaseCoin: goalCoin, QuotedCoin: currentCoin}, 1)
	}
	var quote Quote = Quote{}
	quote.askPair = askPair
	quote.bidPair = bidPair
	quote.goalCoin = goalCoin
	quote.currentCoin = currentCoin
	askVolumeStr := fmt.Sprintf("%g", askPair.Volume)

	// 找出小數點後有幾位
	array := strings.Split(askVolumeStr, ".")
	lastItem := ""
	if strings.Index(askVolumeStr, ".") > 0 {
		lastItem = array[len(array)-1]
	}

	quote.underDot = len(lastItem)
	if val, ok := marketMap[marketName]; ok {
		quote.underDot = val
	}

	return quote
}

func strToFloat64(str string, len int) float64 {
	lenstr := "%." + strconv.Itoa(len) + "f"
	value, _ := strconv.ParseFloat(str, 64)
	value = math.Floor(math.Pow10(len)*value) / math.Pow10(len) // 無條件捨去
	nstr := fmt.Sprintf(lenstr, value)
	val, _ := strconv.ParseFloat(nstr, 64)
	return val
}

type DealFlow struct {
	quotes   []Quote
	takerFee float64
}

func (tri *Triangular) NewDealFlow(goalCoin string, stepCoins []string) DealFlow {
	dealFlow := DealFlow{}
	dealFlow.takerFee = tri.exchangeClient.GetFee().Taker

	tempGoalCoin := goalCoin
	stepLen := len(stepCoins)
	dealFlow.quotes = make([]Quote, stepLen)
	for i := 0; i < stepLen; i++ {
		stepCoin := stepCoins[i]
		quote := tri.NewQuote(tempGoalCoin, stepCoin)
		dealFlow.quotes[i] = quote
		tempGoalCoin = stepCoin
	}
	return dealFlow
}

func (df *DealFlow) getName() string {
	marketName := ""
	for _, value := range df.quotes {
		marketName = fmt.Sprintf("%s%s", marketName, value.MarketName())
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

func (tri *Triangular) stratStrategy() int {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) // 這已經是頂層的 UI 介面了，想以自己的方式呈現錯誤
		}
		mFailCount = mFailCount + 1
		if mFailCount < MAX_FAIL_COUNT {
			//再重來一次
			tri.stratStrategy()
		}
		// 失敗次數太多，直接結束
	}()

	dealFlows := []DealFlow{}

	for _, coinStrip := range tri.CoinStrip {
		fuDealFlow := tri.NewDealFlow(coinStrip.Coins[0], coinStrip.Coins[1:])
		dealFlows = append(dealFlows, fuDealFlow)
	}

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

	currentOrderTotalValue := RANK_N[R_TOTAL_VALUE]

	// 表示有人來搶單拉!!
	if m_expectedTotalValue != 0 && m_expectedTotalValue > currentOrderTotalValue {
		m_isFullPower = true

	} else if profit > RANK_S[R_PROFIT] {
		// 利潤超高 買起來!!!
		m_isFullPower = true
	} else if m_expectedLowestProfit != 0 && m_expectedLowestProfit < profit {
		// 利潤變少了，全力買起來
		m_isFullPower = true
	}

	if m_isFullPower {
		currentOrderTotalValue = RANK_S[R_TOTAL_VALUE]
	}

	wnatOrderTotalValue := currentOrderTotalValue * (1 - (10 * rand.Float64() / 100.0)) // 隨機 -10%
	wnatOrderTotalValue = math.Floor(wnatOrderTotalValue)

	orderTotalValue := math.Min(wnatOrderTotalValue, minSourceTotalValue)

	m_expectedTotalValue = minSourceTotalValue - orderTotalValue
	m_expectedLowestProfit = profit

	// 有利可圖
	if !canOrder(profit, orderTotalValue) {
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
			tri.executeOrder(lowestAskFlow, exc.Ask, laOrderVolume)
			tri.executeOrder(highestBidFlow, exc.Bid, hbOrderVolume)
		} else {
			tri.executeOrder(highestBidFlow, exc.Bid, laOrderVolume)
			tri.executeOrder(lowestAskFlow, exc.Ask, hbOrderVolume)
		}
	}

	content := fmt.Sprintf("%s\r\n %s,\r\n orderTotalValue:%g \r\n profit:%g \r\n m_expectedTotalValue:%g",
		fmt.Sprintf("resAsk:%f, orderVolume:%f, AskCoin:%s", laPrice, laOrderVolume, laName),
		fmt.Sprintf("resBid:%f, orderVolume:%f, bidCoin:%s", hbPrice, hbOrderVolume, hbName),
		orderTotalValue,
		profit,
		m_expectedTotalValue)
	message_tool.SendTelegram(content)

	resPlusSecond := RANK_N[R_PLUS_SECOND]
	if m_isFullPower {
		resPlusSecond = RANK_S[R_PLUS_SECOND]
	}

	return int(resPlusSecond)
}

func canOrder(profit, orderTotalValue float64) bool {
	// 有利可圖
	if profit < 0 {
		log.Println("No profit")
		return false
	} else if profit < RANK_N[R_PROFIT] {
		log.Println("No enough profit")
		return false
	} else if orderTotalValue < mLeastTotalValue {
		log.Println(fmt.Sprintf("orderTotalValue < %f", mLeastTotalValue))
		return false
	}
	return true
}

///

func (tri *Triangular) executeOrder(df DealFlow, pType exc.PriceType, startVolume float64) {
	log.Println(fmt.Sprintf("startVolume:%f", startVolume))
	side := ""
	orderSymbol := 1.0
	switch pType {
	case exc.Bid:
		side = "sell"
		orderVolume := startVolume
		orderSymbol = -1
		for _, quote := range df.quotes {
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*mRangePremium)

			var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
				Market:    quote.MarketName(),
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: exc.MARKET,
			}
			tri.PostOrderRefry(myOrder)

			orderVolume = orderVolume * quote.GetPair(pType).Price
		}
	case exc.Ask:
		side := "buy"
		orderVolume := startVolume
		orderSymbol = 1.0
		var orders []exc.ExchangeOrder = []exc.ExchangeOrder{}
		for i := 0; i < len(df.quotes); i++ {
			quote := df.quotes[i]
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*mRangePremium)

			var myOrder exc.ExchangeOrder = exc.ExchangeOrder{
				Market:    quote.MarketName(),
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: exc.MARKET,
			}
			orders = append(orders, myOrder)

			orderVolume = orderVolume * quote.GetPair(pType).Price

		}

		for i := len(orders) - 1; i >= 0; i-- {
			order := orders[i]
			tri.PostOrderRefry(order)
		}
	}
}

func (tri *Triangular) PostOrderRefry(order exc.ExchangeOrder) {
	//最多重試三次
	const RETRY_TIMES = 3
	for i := 0; i < RETRY_TIMES; i++ {
		_, err := tri.exchangeClient.PostOrder(order)
		if err == nil {
			break
		}
		log.Println(fmt.Sprintf("retry Count:%d", i))
	}

}
