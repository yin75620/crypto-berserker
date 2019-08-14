package main

import (
	"fmt"
	"log"
	"math"
	"net/http"

	fx "github.com/yin75620/crypto-berserker/ftx"
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
	checkProfit()
	stratStrategy()
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
	res := ftx.GetOrderBookResponse(marketName, 1)
	askPair, _ := res.Result.GetPair(1, fx.Ask)
	bidPair, _ := res.Result.GetPair(1, fx.Bid)
	var quote Quote = Quote{}
	quote.askPair = askPair
	quote.bidPair = bidPair
	quote.goalCoin = goalCoin
	quote.currentCoin = currentCoin
	return quote
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
		marketName = fmt.Sprintf("%s%s", marketName, value.MarketName)
	}
	return marketName
}

func (df *DealFlow) getFinalPair(pType fx.PriceType) fx.PricePair {
	return df.getFinalPairWithFee(pType, true)
}

func (df *DealFlow) getFinalPairWithFee(pType fx.PriceType, hasFee bool) fx.PricePair {
	var finalPrice float64 = 1
	var finalVolume float64 = math.MaxFloat64
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
		finalVolume = math.Min(pair.Price*pair.Volume, finalVolume)
	}

	var finalAskPair fx.PricePair = fx.PricePair{}
	finalAskPair.Price = finalPrice
	finalAskPair.Volume = finalVolume
	return finalAskPair
}

func profitCheck(quotes []Quote) {

}

func stratStrategy() {
	fuDealFlow := NewDealFlow(FTT, USD)
	fbuDealFlow := NewDealFlow(FTT, BTC, USD)

	dealFlows := []DealFlow{fuDealFlow, fbuDealFlow}
	lowestAskFlow := getLowestFlow(dealFlows, fx.Ask)
	highestBidFlow := getHighestFlow(dealFlows, fx.Bid)
	lName := lowestAskFlow.getName()
	hName := highestBidFlow.getName()
	fmt.Println(fmt.Sprintf("resAsk:%f, AskCoin:%s", lowestAskFlow.getFinalPair(fx.Ask).Price, lName))
	fmt.Println(fmt.Sprintf("resBid:%f, bidCoin:%s", highestBidFlow.getFinalPair(fx.Bid).Price, hName))

	//fuAskPair := fuDealFlow.getFinalPair(fx.Ask)
	//fbuAskPair := fbuDealFlow.getFinalPair(fx.Ask)

	//pairs := []fx.PricePair{fuAskPair, fbuAskPair}
	//fmt.Println(askLowest)

}

func getLowestFlow(dealFlows []DealFlow, pType fx.PriceType) DealFlow {
	lowest := math.MaxFloat64
	resDealFlow := DealFlow{}
	for _, value := range dealFlows {
		pair := value.getFinalPair(pType)
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
		if highest < pair.Price {
			highest = pair.Price
			resDealFlow = value
		}
	}
	return resDealFlow
}

/*
func startDeal() {

	// 獲取資料
	var fuQuote Quote = NewQuoto("FTT/USD")
	var fbQuote Quote = NewQuoto("FTT/BTC")
	var buQuote Quote = NewQuoto("BTC/USD")

	// 檢查是否有利潤
	//takerFee := 1 + TAKER_FEE
	fbuPrice := fuQuote.price * buQuote.price
	fbuVolume := math.Min(fbQuote.volume*fbQuote.price, buQuote.volume)

	var fbuPair fx.PricePair = fx.PricePair{
		Price:  fbuPrice,
		Volume: fbuVolume,
	}

	coinPairs := []fx.PricePair{fuPair, fbuPair}

	var resPair fx.PricePair
	var resAsk float64 = math.MaxFloat64
	for i := 0; i < len(coinPairs); i++ {
		currentPair := coinPairs[i]

		fmt.Println(fmt.Sprintf("currentAsk:%f, currentCoin:%s", currentPair.Price, currentPair.Volume))
		if currentPair.Price < resAsk {
			resAsk = currentPair.Price
			resPair = currentPair
		}
	}

	fmt.Println(resPair)

	// 檢查是否有利潤

	// 執行套利
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
*/
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
