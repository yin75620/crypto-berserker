package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	gomail "github.com/alexcesaro/mail/gomail"
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
	TAKER_FEE = 0.000665
)

var ftx = fx.NewFtx(http.DefaultClient)

func main() {
	//checkProfit()
	stratStrategy()
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		stratStrategy()
	}

}

func testOrder() {
	marketName := "FTT/USD"
	var myOrder fx.FtxOrder = fx.FtxOrder{
		Market: marketName,
		Side:   "sell",
		Price:  1.82,
		Size:   1,
	}
	ftx.PostOrder(myOrder)
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
	s := fmt.Sprintf("%g", askPair.Volume)
	quote.underDot = strings.LastIndex(s, ".") + 1
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
	var finalVolume float64 = math.MaxFloat64
	var compareVolume float64 = math.MaxFloat64
	for _, quote := range df.quotes {
		pair := quote.GetPair(pType)
		finalPrice = finalPrice * pair.Price
		if hasFee {
			cacularSymbol := 1.0
			if pType == fx.Bid {
				cacularSymbol = -1
			}
			finalPrice = finalPrice * (1.0 + TAKER_FEE*cacularSymbol)
		}
		compareVolume = math.Min(pair.Volume, finalVolume)
		finalVolume = pair.Price * compareVolume
	}

	var finalAskPair fx.PricePair = fx.PricePair{}
	finalAskPair.Price = finalPrice
	finalAskPair.Volume = finalVolume / finalPrice
	return finalAskPair
}
func getLowestFlow(dealFlows []DealFlow, pType fx.PriceType) DealFlow {
	lowest := math.MaxFloat64
	resDealFlow := DealFlow{}
	for _, value := range dealFlows {
		pair := value.getFinalPair(pType)
		fmt.Println(fmt.Sprintf("getLowestFlow:%f, Coin:%s", pair.Price, value.getName()))
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
		fmt.Println(fmt.Sprintf("getHighestFlow:%f, Coin:%s", pair.Price, value.getName()))
		if highest < pair.Price {
			highest = pair.Price
			resDealFlow = value
		}
	}
	return resDealFlow
}

func stratStrategy() {
	fuDealFlow := NewDealFlow(FTT, USD)
	fbuDealFlow := NewDealFlow(FTT, BTC, USD)
	futDealFlow := NewDealFlow(FTT, USDT, USD)

	dealFlows := []DealFlow{fuDealFlow, fbuDealFlow, futDealFlow}
	lowestAskFlow := getLowestFlow(dealFlows, fx.Ask)
	highestBidFlow := getHighestFlow(dealFlows, fx.Bid)
	laName := lowestAskFlow.getName()
	hbName := highestBidFlow.getName()

	laPrice := lowestAskFlow.getFinalPair(fx.Ask).Price
	hbPrice := highestBidFlow.getFinalPair(fx.Bid).Price

	laVolume := lowestAskFlow.getFinalPair(fx.Ask).Volume
	hbVolume := highestBidFlow.getFinalPair(fx.Bid).Volume

	orderVolume := math.Min(laVolume, hbVolume)

	const PER_ORDER_MAX_VOLUME = 200
	orderVolume = math.Min(PER_ORDER_MAX_VOLUME, orderVolume)

	fmt.Println(fmt.Sprintf("resAsk:%f, AskCoin:%s", laPrice, laName))
	fmt.Println(fmt.Sprintf("resBid:%f, bidCoin:%s", hbPrice, hbName))

	profit := hbPrice - laPrice
	fmt.Println(fmt.Sprintf("Profit:%f", profit))

	// 有利可圖
	if profit < 0 {
		fmt.Println("No profit")
		return
	} else if profit < 0.0001 {
		fmt.Println("No enough profit")
		return
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

	content := fmt.Sprintf("%s%s, volume:%g", fmt.Sprintf("resAsk:%f, AskCoin:%s", laPrice, laName), fmt.Sprintf("resBid:%f, bidCoin:%s", hbPrice, hbName), orderVolume)
	sendMail(content)
}

/// email
func sendMail(content string) {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "yin75620@gmail.com", "Golang")
	msg.SetHeader("To", "yin75620@gmail.com")
	msg.AddHeader("To", "yin75620@gmail.com")
	msg.SetHeader("Subject", "Hello!")
	msg.SetBody("text/plain", "Hello Has Profit")
	msg.AddAlternative("text/html", "content")

	m := gomail.NewMailer("smtp.gmail.com", "yin75620", setting.GMAIL_PASSWORD, 25)
	if err := m.Send(msg); err != nil {
		log.Println(err)
	}
}

///

func executeOrder(df DealFlow, pType fx.PriceType, startVolume float64) {
	fmt.Println(fmt.Sprintf("startVolume:%f", startVolume))
	side := ""
	switch pType {
	case fx.Bid:
		side = "sell"
		orderVolume := startVolume
		for _, quote := range df.quotes {
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price
			postOrder(quote.MarketName(), side, orderPrice, orderVolume)
			orderVolume = orderVolume * quote.GetPair(pType).Price
		}
		/*case fx.Ask:
		side := "buy"
		orderVolume := startVolume
		for i := len(df.quotes) - 1; i >= 0; i-- {
			quote := df.quotes[i]
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)
			orderPrice := quote.GetPair(pType).Price
			postOrder(quote.MarketName(), side, orderPrice, orderVolume)
			orderVolume = orderVolume * quote.GetPair(pType).Price
		}*/
		/*
			case fx.Ask:
				side := "buy"
				orderVolume := startVolume
				var orders []fx.FtxOrder = []fx.FtxOrder{}
				for i := len(df.quotes) - 1; i >= 0; i-- {
					quote := df.quotes[i]
					orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

					orderPrice := quote.GetPair(pType).Price

					var myOrder fx.FtxOrder = fx.FtxOrder{
						Market: quote.MarketName(),
						Side:   side,
						Price:  orderPrice,
						Size:   orderVolume,
					}
					orders = append(orders, myOrder)

					orderVolume = orderVolume * quote.GetPair(pType).Price

				}

				for i := 0; i < len(orders); i++ {
					order := orders[i]
					ftx.PostOrder(order)
				}
			}*/
	case fx.Ask:
		side := "buy"
		orderVolume := startVolume
		var orders []fx.FtxOrder = []fx.FtxOrder{}
		for i := 0; i < len(df.quotes); i++ {
			quote := df.quotes[i]
			orderVolume = strToFloat64(fmt.Sprintf("%g", orderVolume), quote.underDot)

			orderPrice := quote.GetPair(pType).Price

			var myOrder fx.FtxOrder = fx.FtxOrder{
				Market: quote.MarketName(),
				Side:   side,
				Price:  orderPrice,
				Size:   orderVolume,
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

func checkProfit() {
	resAsk, AskCoin := getLowestAsk([]string{USD, BTC}, FTT, USD)

	resBid, bidCoin := getTopBid([]string{USD, BTC}, FTT, USD)
	fmt.Println(fmt.Sprintf("resAsk:%f, AskCoin:%s", resAsk, AskCoin))
	fmt.Println(fmt.Sprintf("resBid:%f, bidCoin:%s", resBid, bidCoin))

	// 檢查是否有利潤
	profit := (resBid - resAsk) / resAsk
	fmt.Println(fmt.Sprintf("profit:%f", profit))

	// 有利潤
	if profit < 0 {
		// 執行套利
		fmt.Println("No Profit")
		return
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

///
//
// currentCoin: 目前幣別
// goalCoin: 兌換目標幣別
// baseCoin: 計價幣別
// 回傳計價幣別 數量
func getAskWithBaseCoin(goalCoin, currentCoin, baseCoin string) float64 {

	if currentCoin == baseCoin {
		return getAsk(goalCoin, currentCoin)
	}

	gcAsk := getAsk(goalCoin, currentCoin)
	cbAsk := getAsk(currentCoin, baseCoin)
	res := gcAsk * cbAsk * (1 + TAKER_FEE)
	return res
}

func getAsk(goalCoin, currentCoin string) float64 {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)
	res := ftx.GetAsk(marketName, 1)
	//fmt.Printf("%f", res)

	return res
}

func getBidWithBaseCoin(goalCoin, currentCoin, baseCoin string) float64 {

	if currentCoin == baseCoin {
		return getBid(goalCoin, currentCoin)
	}

	gcBid := getBid(goalCoin, currentCoin)
	cbBid := getBid(currentCoin, baseCoin)
	res := gcBid * cbBid * (1 - TAKER_FEE)
	return res
}

func getBid(goalCoin, currentCoin string) float64 {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)
	res := ftx.GetBid(marketName, 1)
	//fmt.Printf("%f", res)

	return res
}

// 最高收購價
func getTopBid(coins []string, goalCoin string, baseCoin string) (float64, string) {

	//價格比較
	resCoin := ""
	var resBid float64 = 0
	for i := 0; i < len(coins); i++ {
		coin := coins[i]
		currentBid := getBidWithBaseCoin(goalCoin, coin, baseCoin)

		fmt.Println(fmt.Sprintf("currentBid:%f, currentCoin:%s", currentBid, coin))
		if currentBid > resBid {
			resBid = currentBid
			resCoin = coin
		}
	}

	return resBid, resCoin
}

// 最低求售價
func getLowestAsk(coins []string, goalCoin string, baseCoin string) (float64, string) {

	//交易所查詢，先想像只有這個，多個交易所就多個互相比較後就知道了

	//價格比較
	resCoin := ""
	var resAsk float64 = math.MaxFloat64
	for i := 0; i < len(coins); i++ {
		coin := coins[i]
		currentAsk := getAskWithBaseCoin(goalCoin, coin, baseCoin)

		fmt.Println(fmt.Sprintf("currentAsk:%f, currentCoin:%s", currentAsk, coin))
		if currentAsk < resAsk {
			resAsk = currentAsk
			resCoin = coin
		}
	}

	return resAsk, resCoin
}
