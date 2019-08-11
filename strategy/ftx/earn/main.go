package main

import (
	"fmt"
	"math"
	"net/http"

	fx "github.com/yin75620/crypto-berserker/ftx"
)

const (
	USD  = "USD"
	USDT = "USDt"
	BTC  = "BTC"
	FTT  = "FTT"
)

const (
	TAKER_FEE = 0.000665
)

var ftx = fx.NewFtx(http.DefaultClient)

func main() {
	checkProfit()
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
	res := gcAsk * cbAsk
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
	res := gcBid * cbBid
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
		if coin != baseCoin {
			currentAsk = currentAsk * (1 + TAKER_FEE)
		}
		fmt.Println(fmt.Sprintf("currentAsk:%f, currentCoin:%s", currentAsk, coin))
		if currentAsk < resAsk {
			resAsk = currentAsk
			resCoin = coin
		}
	}

	return resAsk, resCoin
}

type CoinEchange struct {
}
