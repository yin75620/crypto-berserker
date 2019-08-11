package main

import (
	"fmt"
	"log"
	"net/http"

	fx "github.com/yin75620/crypto-berserker/ftx"
)

const (
	USD  = "USD"
	USDT = "USDt"
	BTC  = "BTC"
	FTT  = "FTT"
)

func main() {
	// USD -> FTT

	res := getAsk(FTT, USD)
	log.Println(res)
	// USD/BTC
	// USD/FTT
	// BTC/FTT

	// 檢查是否有利潤

	// 執行套利
	//
}

///
// baseCoin: 計價幣別
// currentCoin: 目前幣別
// goalCoin: 兌換目標幣別
// 回傳計價幣別 數量
func getPrize(baseCoin, goalCoin, currentCoin string) float64 {
	// 查 base/goal
	//getPrize()
	return 0
}

func getAsk(goalCoin, currentCoin string) float64 {
	marketName := fmt.Sprintf("%s/%s", goalCoin, currentCoin)

	var ftx = fx.NewFtx(http.DefaultClient)

	res := ftx.GetAsk(marketName, 1)
	fmt.Printf("%f", res)

	return 0
}

// 最高收購價
func getTopBid() {

}

// 最低求售價
func getLowestAsk(coins []string, goalCoin string) (float64, string) {
	return 0.01, ""
	//交易所查詢，先想像只有這個，多個交易所就多個互相比較後就知道了

	//價格比較

}

type CoinEchange struct {
}
