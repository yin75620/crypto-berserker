package main

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

	gomail "github.com/alexcesaro/mail/gomail"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	fx "github.com/yin75620/crypto-berserker/ftx"
	"github.com/yin75620/crypto-berserker/setting"
)

const (
	USD  = "USD"
	USDT = "USDT"
	BTC  = "BTC"
	FTT  = "FTT"
)

const (
	TAKER_FEE            = 0.000665
	RANGE_PREMIUM        = 0.2 //20%
	PER_ORDER_MAX_VOLUME = 853 //有人搶就全力對搶
	PROFIT_THRESHOLD     = 0.001
	LEAST_VOLUME         = 10
)

//當const 用
var (
	// 數量, 利潤, 加速
	RANK_S = []float64{PER_ORDER_MAX_VOLUME, 0.01, -2.0}
	RANK_N = []float64{PER_ORDER_MAX_VOLUME / 4, 0.001, 0.0}
)

const (
	R_VOLUME      = 0
	R_PROFIT      = 1
	R_PLUS_SECOND = 2
)

var ftx = fx.NewFtx(http.DefaultClient)

func main() {
	StartTelegram()
	var logFile *os.File = StartLog()
	defer logFile.Close()
	stratStrategy()

	infoStr := string(ftx.GetAccountInfo())
	sendTelegram(infoStr)
	/*
		ticker := time.NewTicker(6 * time.Second)
		for _ = range ticker.C {
			stratStrategy()
		}*/
	var delay_time int = 5
	d := time.Duration(time.Second * time.Duration(delay_time))

	t := time.NewTimer(d)
	defer t.Stop()

	for {
		<-t.C
		plusSecond := stratStrategy()
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
	askPair     fx.PricePair
	bidPair     fx.PricePair
	underDot    int
}

func (q *Quote) MarketName() string {
	marketName := fmt.Sprintf("%s/%s", q.goalCoin, q.currentCoin)
	return marketName
}

func (q *Quote) GetPair(pType fx.PriceType) fx.PricePair {
	switch pType {
	case fx.Ask:
		return q.askPair
	case fx.Bid:
		return q.bidPair
	}
	log.Fatal("Error: not specific pType")
	return fx.PricePair{}
}

func NewQuote(goalCoin, currentCoin string) Quote {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)
	var askPair fx.PricePair
	var bidPair fx.PricePair
	//偽裝成 USDT
	if marketName == "USDT/USD" {
		askPair = fx.PricePair{1.001, 999999}
		bidPair = fx.PricePair{0.997, 999999}
	} else {
		res := ftx.GetOrderBookResponse(marketName, 1)
		askPair, _ = res.Result.GetPair(1, fx.Ask)
		bidPair, _ = res.Result.GetPair(1, fx.Bid)
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
	quotes []Quote
}

func NewDealFlow(goalCoin string, stepCoins ...string) DealFlow {
	var dealFlow DealFlow = DealFlow{}
	tempGoalCoin := goalCoin
	stepLen := len(stepCoins)
	dealFlow.quotes = make([]Quote, stepLen)
	for i := 0; i < stepLen; i++ {
		stepCoin := stepCoins[i]
		quote := NewQuote(tempGoalCoin, stepCoin)
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

func (df *DealFlow) getFinalPair(pType fx.PriceType) fx.PricePair {
	return df.getFinalPairWithFee(pType, true)
}

func (df *DealFlow) getFinalPairWithFee(pType fx.PriceType, hasFee bool) fx.PricePair {
	var finalPrice float64 = 1
	var finalFeePrice float64 = 1
	var finalVolume float64 = math.MaxFloat64
	var compareVolume float64 = math.MaxFloat64
	for _, quote := range df.quotes {
		pair := quote.GetPair(pType)
		finalPrice = finalPrice * pair.Price
		finalFeePrice = finalFeePrice * pair.Price
		if hasFee {
			cacularSymbol := 1.0
			if pType == fx.Bid {
				cacularSymbol = -1
			}
			finalFeePrice = finalFeePrice * (1.0 + TAKER_FEE*cacularSymbol)
		}
		compareVolume = math.Min(pair.Volume, finalVolume)
		finalVolume = pair.Price * compareVolume
	}

	var finalAskPair fx.PricePair = fx.PricePair{}
	finalAskPair.Price = finalFeePrice
	finalAskPair.Volume = finalVolume / finalPrice
	return finalAskPair
}
func getLowestFlow(dealFlows []DealFlow, pType fx.PriceType) DealFlow {
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

func getHighestFlow(dealFlows []DealFlow, pType fx.PriceType) DealFlow {
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
	m_expectedSourceOrder  float64 = 0
	m_expectedLowestProfit float64 = 0
)

var m_isFullPower = false

func stratStrategy() int {
	fuDealFlow := NewDealFlow(FTT, USD)
	fbuDealFlow := NewDealFlow(FTT, BTC, USD)
	//futDealFlow := NewDealFlow(FTT, USDT, USD)

	dealFlows := []DealFlow{
		fuDealFlow,
		fbuDealFlow,
		//futDealFlow,
	}
	lowestAskFlow := getLowestFlow(dealFlows, fx.Ask)
	highestBidFlow := getHighestFlow(dealFlows, fx.Bid)
	laName := lowestAskFlow.getName()
	hbName := highestBidFlow.getName()

	laPrice := lowestAskFlow.getFinalPair(fx.Ask).Price
	hbPrice := highestBidFlow.getFinalPair(fx.Bid).Price

	laVolume := lowestAskFlow.getFinalPair(fx.Ask).Volume
	hbVolume := highestBidFlow.getFinalPair(fx.Bid).Volume

	sourceOrderVolume := math.Min(laVolume, hbVolume)

	log.Println(fmt.Sprintf("sourceOrderVolume:%g", sourceOrderVolume))
	log.Println(fmt.Sprintf("m_expectedSourceOrder:%g", m_expectedSourceOrder))

	log.Println(fmt.Sprintf("resAsk:%f, AskCoin:%s", laPrice, laName))
	log.Println(fmt.Sprintf("resBid:%f, bidCoin:%s", hbPrice, hbName))

	profit := hbPrice - laPrice

	log.Println(fmt.Sprintf("Profit:%f", profit))

	perOrderMaxVolume := RANK_N[R_VOLUME]

	// 表示有人來搶單拉!!
	if m_expectedSourceOrder != 0 && m_expectedSourceOrder > sourceOrderVolume {
		m_isFullPower = true

	} else if profit > RANK_S[R_PROFIT] {
		// 利潤超高 買起來!!!
		m_isFullPower = true
	} else if m_expectedLowestProfit != 0 && m_expectedLowestProfit < profit {
		// 利潤變少了，全力買起來
		m_isFullPower = true
	}

	if m_isFullPower {
		perOrderMaxVolume = RANK_S[R_VOLUME]
	}

	wnatOrderVolume := perOrderMaxVolume * (1 + (10 * rand.Float64() / 100.0)) // 隨機 +10%
	wnatOrderVolume = math.Floor(wnatOrderVolume)

	orderVolume := math.Min(wnatOrderVolume, sourceOrderVolume)

	m_expectedSourceOrder = sourceOrderVolume - orderVolume
	m_expectedLowestProfit = profit

	// 有利可圖
	if !canOrder(profit, orderVolume) {
		// 無利可圖，重設偵測
		m_isFullPower = false
		m_expectedSourceOrder = 0
		m_expectedLowestProfit = 0
		return 0
	}

	const isOrder = true
	if isOrder {
		const isKeepUSD = true
		if isKeepUSD {
			executeOrder(lowestAskFlow, fx.Ask, orderVolume)
			executeOrder(highestBidFlow, fx.Bid, orderVolume)
		} else {
			executeOrder(highestBidFlow, fx.Bid, orderVolume)
			executeOrder(lowestAskFlow, fx.Ask, orderVolume)
		}
	}

	content := fmt.Sprintf("%s\r\n %s,\r\n volume:%g \r\n profit:%g \r\n sourceVolume:%g",
		fmt.Sprintf("resAsk:%f, AskCoin:%s", laPrice, laName),
		fmt.Sprintf("resBid:%f, bidCoin:%s", hbPrice, hbName),
		orderVolume,
		profit,
		sourceOrderVolume)
	sendTelegram(content)

	resPlusSecond := RANK_N[R_PLUS_SECOND]
	if m_isFullPower {
		resPlusSecond = RANK_S[R_PLUS_SECOND]
	}

	return int(resPlusSecond)
}

func canOrder(profit, orderVolume float64) bool {
	// 有利可圖
	if profit < 0 {
		log.Println("No profit")
		return false
	} else if profit < PROFIT_THRESHOLD {
		log.Println("No enough profit")
		return false
	} else if orderVolume < LEAST_VOLUME {
		log.Println("orderVolume < 10")
		return false
	}
	return true
}

/// message
func sendMail(content string) {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "yin75620@gmail.com", "Golang")
	msg.SetHeader("To", "yin75620@gmail.com")
	msg.AddHeader("To", "yin75620@gmail.com")
	msg.SetHeader("Subject", "Hello!")
	msg.SetBody("text/plain", "Hello Has Profit")
	msg.AddAlternative("text/html", content)

	m := gomail.NewMailer("smtp.gmail.com", "yin75620", setting.GMAIL_PASSWORD, 25)
	if err := m.Send(msg); err != nil {
		log.Println(err)
	}
}

var bot *tgbotapi.BotAPI

func StartTelegram() {
	bot, _ = tgbotapi.NewBotAPI(setting.TELEGRAM_BOT_TOKEN)
}

func sendTelegram(content string) {
	msg := tgbotapi.NewMessage(945156610, content)

	bot.Send(msg)
}

///

func executeOrder(df DealFlow, pType fx.PriceType, startVolume float64) {
	log.Println(fmt.Sprintf("startVolume:%f", startVolume))
	side := ""
	orderSymbol := 1.0
	switch pType {
	case fx.Bid:
		side = "sell"
		orderVolume := startVolume
		orderSymbol = -1
		for _, quote := range df.quotes {
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*RANGE_PREMIUM)

			var myOrder fx.FtxOrder = fx.FtxOrder{
				Market:    quote.MarketName(),
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: fx.MARKET,
			}
			ftx.PostOrder(myOrder)

			orderVolume = orderVolume * quote.GetPair(pType).Price
		}
	case fx.Ask:
		side := "buy"
		orderVolume := startVolume
		orderSymbol = 1.0
		var orders []fx.FtxOrder = []fx.FtxOrder{}
		for i := 0; i < len(df.quotes); i++ {
			quote := df.quotes[i]
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price

			myOrderPrice := orderPrice * (1 + orderSymbol*RANGE_PREMIUM)

			var myOrder fx.FtxOrder = fx.FtxOrder{
				Market:    quote.MarketName(),
				Side:      side,
				Price:     myOrderPrice,
				Size:      orderVolume,
				OrderType: fx.MARKET,
			}
			orders = append(orders, myOrder)

			orderVolume = orderVolume * quote.GetPair(pType).Price

		}

		for i := len(orders) - 1; i >= 0; i-- {
			order := orders[i]
			ftx.PostOrder(order)
		}
	}
}

///
//交易
func postCoinOrder(goalCoin, currentCoin, side string, price, size float64) {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)
	postOrder(marketName, side, price, size)
}

func postOrder(marketName, side string, price, size float64) {
	var myOrder fx.FtxOrder = fx.FtxOrder{
		Market: marketName,
		Side:   side,
		Price:  price,
		Size:   size,
	}
	ftx.PostOrder(myOrder)
}
