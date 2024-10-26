package main

import (
	"fmt"
	"time"

	"github.com/yin75620/crypto-berserker/exchange-list/common"
)

func main() {
	collector := common.GetCollector("")

	// startTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	// endTime, _ := time.Parse(time.RFC3339, "2024-01-02T00:00:00Z")
	startTime := time.Time{}
	endTime := time.Time{}

	klines, err := collector.GetKlines("BTCUSDT", "1m", startTime, endTime, 3)
	if err != nil {
		fmt.Println("Error fetching Klines:", err)
		return
	}

	// 顯示取得的 K 線資料
	for _, kline := range klines {
		fmt.Printf("Open: %f, Close: %f, Volume: %f\n", kline.Open, kline.Close, kline.Volume)
	}
}
